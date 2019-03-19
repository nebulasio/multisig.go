"use strict";

let crypto = require('crypto.js');

function MultiSign() {
    this.decimal = new BigNumber('1000000000000000000');
    this.managers = MANAGERS;
    this._managers = null;
    this._deletedManagers = null;
    this._addedManagers = null;
    this._sysConfig = null;
    this._sendNasRule = null;

    this.dataKeyPrefix = "data_";
    this.deletedManagersKey = "deleted_managers";
    this.addedManagersKey = "added_managers";
    this.managerUpdateLogKey = "manager_update_log";
    this.sendNasRuleKey = "send_nas_rule";
    this.sendNasRuleUpdateLogKey = "send_nas_rule_update_log";
    this.sysConfigKey = "sys_config";
    this.sysConfigUpdateLogKey = "sys_config_update_log";

    this.actionDeleteManager = "delete_manager";
    this.actionAddManager = "add_manager";
    this.actionReplaceManager = "replace_manager";
    this.actionSendNas = "send_nas";
    this.actionUpdateSysConfig = "update_sys_config";
    this.actionUpdateSendNasRule = "update_send_nas_rule";

    this.infinity = "INFINITY";

    LocalContractStorage.defineMapProperty(this, "data", {
        parse: function (text) {
            return JSON.parse(text);
        },
        stringify: function (obj) {
            return JSON.stringify(obj);
        }
    });
}

MultiSign.prototype = {

    init: function () {
    },

    _setData: function (key, data) {
        this.data.put(this.dataKeyPrefix + key, data);
    },

    _getData: function (key) {
        return this.data.get(this.dataKeyPrefix + key);
    },

    _getSysConfig: function () {
        if (!this._sysConfig) {
            this._sysConfig = this.data.get(this.sysConfigKey);
        }
        if (!this._sysConfig) {
            this._sysConfig = {
                "version": "0",
                "proportionOfSigners": {
                    "updateSysConfig": "1",
                    "updateSendNasRule": "1",
                    "addManager": "1",
                    "deleteManager": "1",
                    "replaceManager": "1"
                }
            }
        }
        return this._sysConfig;
    },

    _getSendNasRule: function () {
        if (!this._sendNasRule) {
            this._sendNasRule = this.data.get(this.sendNasRuleKey);
        }
        if (!this._sendNasRule) {
            // TODO: Confirm the default rule when going online
            this._sendNasRule = {
                "version": "0",
                "rule": [
                    {
                        "startValue": "0",
                        "endValue": "0.5",
                        "proportionOfSigners": "0.3"
                    },
                    {
                        "startValue": "0.5",
                        "endValue": "1",
                        "proportionOfSigners": "0.5"
                    },
                    {
                        "startValue": "1",
                        "endValue": this.infinity,
                        "proportionOfSigners": "1"
                    }
                ]
            };
        }
        return this._sendNasRule;
    },

    _getLogs(logKey) {
        return this.data.get(logKey);
    },

    _addLog: function (log, logKey) {
        let logs = this.data.get(logKey);
        if (!logs) {
            logs = [];
        }
        log.timestamp = new Date().getTime();
        logs.push(log);
        this.data.put(logKey, logs);
    },

    _getDeletedManagers: function () {
        if (!this._deletedManagers) {
            this._deletedManagers = this.data.get(this.deletedManagersKey);
        }
        if (!this._deletedManagers) {
            this._deletedManagers = [];
        }
        return this._deletedManagers;
    },

    _getAddedManagers: function () {
        if (!this._addedManagers) {
            this._addedManagers = this.data.get(this.addedManagersKey);
        }
        if (!this._addedManagers) {
            this._addedManagers = [];
        }
        return this._addedManagers;
    },

    _getManagers: function () {
        if (this._managers == null) {
            let deletedManagers = this._getDeletedManagers();
            let addedManagers = this._getAddedManagers();

            this._managers = [];
            for (let i = 0; i < this.managers.length; ++i) {
                this._managers.push(this.managers[i]);
            }
            for (let i = 0; i < addedManagers.length; ++i) {
                this._managers.push(addedManagers[i]);
            }
            for (let i = 0; i < deletedManagers.length; ++i) {
                let index = this._managers.indexOf(deletedManagers[i]);
                this._managers.splice(index, 1);
            }
        }
        return this._managers;
    },

    _isDeleted: function (address) {
        let ms = this._getDeletedManagers();
        return ms.indexOf(address) >= 0;
    },

    _verifyAddresses: function (address) {
        for (let i = 0; i < arguments.length; ++i) {
            if (0 === Blockchain.verifyAddress(arguments[i])) {
                throw (address + ' is not a valid nas address.');
            }
        }
    },

    _verifyDeleteAddress: function (address) {
        let managers = this._getManagers();
        if (managers.indexOf(address) < 0) {
            throw (address + ' could not be found.');
        }
    },

    _verifyAddAddress: function (address) {
        let managers = this._getManagers();
        if (managers.indexOf(address) >= 0) {
            throw (address + ' is already a manager.');
        }
        if (this._isDeleted(address)) {
            throw (address + ' has been deleted.');
        }
    },

    _checkNumbers: function () {
        for (let i = 0; i < arguments.length; ++i) {
            let n = arguments[i];
            if (n === null || n === undefined || !/^\d+(\.\d+)?$/.test(n)) {
                throw ('Data error.');
            }
        }
    },

    _checkProportions: function () {
        for (let i = 0; i < arguments.length; ++i) {
            let p = arguments[i];
            this._checkNumbers(p);
            p = parseFloat(p)
            if (p <= 0 || p > 1) {
                throw ('Proportion error');
            }
        }
    },

    _deleteManagerAddress(address) {
        let deletedManagers = this._getDeletedManagers();
        deletedManagers.push(address);
        this.data.put(this.deletedManagersKey, deletedManagers);
    },

    _addManagerAddress(address) {
        let addedManagers = this._getAddedManagers();
        addedManagers.push(address);
        this.data.put(this.addedManagersKey, addedManagers);
    },

    /**
     * { "action": "delete_manager", "detail": "n1xxxxx" }
     */
    _deleteManager: function (data, signers) {
        let address = data.detail;
        this._verifyAddresses(address);
        this._verifyDeleteAddress(address);
        this._deleteManagerAddress(address);
        data.signers = signers;
        this._addLog(data, this.managerUpdateLogKey);
    },

    /**
     * { "action": "add_manager", "detail": "n1xxxxx" }
     */
    _addManager: function (data, signers) {
        let address = data.detail;
        this._verifyAddresses(address);
        this._verifyAddAddress(address);
        this._addManagerAddress(address);
        data.signers = signers;
        this._addLog(data, this.managerUpdateLogKey);
    },

    /**
     * { "action": "replace_manager", "detail": { "oldAddress": "n1xxxxx", "newAddress": "n1yyyy" } }
     */
    _replaceManager: function (data, signers) {
        let oldAddress = data.detail.oldAddress;
        let newAddress = data.detail.newAddress;
        this._verifyAddresses(oldAddress, newAddress);
        if (oldAddress === newAddress) {
            throw ('Old address cannot be the same as new address')
        }
        this._verifyDeleteAddress(oldAddress);
        this._verifyAddAddress(newAddress);
        this._deleteManagerAddress(oldAddress);
        this._addManagerAddress(newAddress);

        data.signers = signers;
        this._addLog(data, this.managerUpdateLogKey);
    },

    /**
     * { "action": "send_nas", "detail": { "id": "xxx", "to": "n1xxxxx", "value": "0.001" } }
     */
    _sendNas: function (data, signers) {
        let tx = data.detail;
        this._verifyAddAddress(tx.to);
        this._checkNumbers(tx.value);
        let value = new BigNumber(tx.value).times(this.decimal);
        let balance = new BigNumber(Blockchain.getAccountState(Blockchain.transaction.to).balance);
        if (balance.comparedTo(value) < 0) {
            throw ('Insufficient balance.');
        }
        if (this._getData(tx.id)) {
            throw (tx.id + ' exists.');
        }
        if (!Blockchain.transfer(tx.to, value)) {
            throw ('Transfer error.');
        }
        data.signers = signers;
        this._setData(tx.id, data);
    },

    /**
     * {
     *     "action": "update_send_nas_rule",
     *     "detail": {
     *          "version": "0",
     *          "rule": [
     *              { "startValue": "0", "endValue": "0.5", "proportionOfSigners": "0.3" },
     *              { "startValue": "0.5", "endValue": "1", "proportionOfSigners": "0.5" },
     *              { "startValue": "1", "endValue": this.infinity, "proportionOfSigners": "1" }
     *          ]
     *     ]
     * }
     */
    _updateSendNasRule: function (data, signers) {
        this._checkNumbers(data.detail.version);

        let ver = parseFloat(data.detail.version);
        if (ver <= parseFloat(this._getSendNasRule().version)) {
            throw ('Version error.');
        }
        let rule = data.detail.rule;
        if (rule.length <= 0) {
            throw ('Rules cannot be empty.');
        }
        let v = new BigNumber("0");
        for (let r in rule) {
            this._checkNumbers(r.startValue);
            if (r.endValue !== this.infinity) {
                this._checkNumbers(r.endValue);
            }
            this._checkProportions(r.proportionOfSigners);
            let start = new BigNumber(r.startValue);
            let end = r.endValue === this.infinity ? null : new BigNumber(r.endValue);
            if (v == null || v.compareTo(start) !== 0 || (end != null && end.compareTo(start) <= 0)) {
                throw ('Rules error.');
            }
            v = end;
        }
        if (v != null) {
            throw ('Rules error.');
        }
        this._sendNasRule = rule;
        this.data.put(this.sendNasRuleKey, rule);
        data.signers = signers;
        this._addLog(data, this.sendNasRuleUpdateLogKey);
    },

    /**
     * {
     *     "action": "update_sys_config",
     *     "detail": {
     *          "version": "0",
     *          "proportionOfSigners": {
     *              "updateSysConfig": "1",
     *              "updateSendNasRule": "1",
     *              "addManager": "1",
     *              "deleteManager": "1",
     *              "replaceManager": "1"
     *          }
     *     }
     * }
     */
    _updateSysConfig: function (data, signers) {
        this._checkNumbers(data.detail.version);
        let ver = parseFloat(data.detail.version);
        if (ver <= parseFloat(this._getSysConfig().version)) {
            throw ('Version error.');
        }
        let ps = data.detail.proportionOfSigners;
        this._checkProportions(
            ps.updateSysConfig, ps.updateSendNasRule,
            ps.addManager, ps.deleteManager, ps.replaceManager
        );
        let config = data.detail;
        this._sysConfig = config;
        this.data.put(this.sysConfigKey, config);
        data.signers = signers;
        this._addLog(data, this.sysConfigUpdateLogKey);
    },

    _getNumOfNeedsSignersWithSendNasConfig: function (data) {
        let v = new BigNumber(data.detail.value);
        let rule = this._getSendNasRule().rule;
        let n = 1;
        for (let i = 0; i < rule.length; ++i) {
            let r = rule[i];
            let start = new BigNumber(r.startValue);
            let end = r.endValue === this.infinity ? null : new BigNumber(r.endValue);
            if (v.compareTo(start) >= 0 && (end == null || end.compareTo(v) > 0)) {
                n = parseFloat(r.proportionOfSigners);
                break;
            }
        }
        return n * this._getManagers();
    },

    _getNumOfNeedsSigners: function (data) {
        let p = this._getSysConfig().proportionOfSigners;
        switch (data.action) {
            case this.actionAddManager:
                return this._getManagers().length * parseFloat(p.addManager);

            case this.actionDeleteManager:
                return (this._getManagers().length - 1) * parseFloat(p.deleteManager);

            case this.actionReplaceManager:
                return (this._getManagers().length - 1) * parseFloat(p.replaceManager);

            case this.actionUpdateSysConfig:
                return this._getManagers().length * parseFloat(p.updateSysConfig);

            case this.actionUpdateSendNasRule:
                return this._getManagers().length * parseFloat(p.updateSendNasRule);

            case this.actionSendNas:
                return this._getNumOfNeedsSignersWithSendNasConfig(data);

            default:
                throw ('Action ' + data.action + ' is not supported.');
        }
    },

    _verifySign: function (item) {
        let managers = this._getManagers();
        let signers = [];

        let hash = crypto.sha3256(item.data);
        item.data = JSON.parse(item.data);
        let n = this._getNumOfNeedsSigners(item.data);
        if (item.sigs == null || item.sigs.length < n) {
            throw ('Minimum ' + n + ' signatures required.');
        }

        for (let j = 0; j < item.sigs.length; ++j) {
            let s = item.sigs[j];
            let a = crypto.recoverAddress(1, hash, s);
            if (a != null && managers.indexOf(a) >= 0) {
                if (signers.indexOf(a) >= 0) {
                    throw ('Signature "' + s + '" is repeated.');
                } else {
                    signers.push(a);
                }
            } else {
                throw ('Signature "' + s + '" error.');
            }
        }

        return signers;
    },

    _execute: function (data, signers) {
        switch (data.action) {
            case this.actionAddManager:
                this._addManager(data, signers);
                break;
            case this.actionDeleteManager:
                this._deleteManager(data, signers);
                break;
            case this.actionReplaceManager:
                this._replaceManager(data, signers);
                break;
            case this.actionSendNas:
                this._sendNas(data, signers);
                break;
            case this.actionUpdateSendNasRule:
                this._updateSendNasRule(data, signers);
                break;
            case this.actionUpdateSysConfig:
                this._updateSysConfig(data, signers);
                break;
            default:
                throw ('Action ' + data.action + ' is not supported.');
        }
    },

    getManagers: function () {
        return this._getManagers();
    },

    getSysConfig: function () {
        return this._getSysConfig()
    },

    getSendNasRule: function () {
        return this._getSendNasRule();
    },

    getManagerUpdateLogs: function () {
        return this._getLogs(this.managerUpdateLogKey);
    },

    getSendNasRuleUpdateLogs: function () {
        return this._getLogs(this.sendNasRuleUpdateLogKey);
    },

    getSysConfigUpdateLogs: function () {
        return this._getLogs(this.sysConfigUpdateLogKey);
    },

    execute: function (array) {
        array = JSON.parse(array);
        if (array == null || array.length === 0) {
            throw('Data is empty');
        }
        for (let i = 0; i < array.length; ++i) {
            let item = array[i];
            let signers = this._verifySign(item);
            this._execute(item.data, signers);
        }
        return "success";
    },

    query: function (dataId) {
        return this._getData(dataId);
    },

    accept: function () {
        Event.Trigger("transfer", {
            Transfer: {
                from: Blockchain.transaction.from,
                to: Blockchain.transaction.to,
                value: Blockchain.transaction.value,
            }
        });
    }
};

module.exports = MultiSign;
