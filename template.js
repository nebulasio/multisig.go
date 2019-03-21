"use strict";

let crypto = require('crypto.js');

function MultiSign() {
    this.decimal = new BigNumber('1000000000000000000');
    this.signees = SIGNEES;
    this._signees = null;
    this._removedSignees = null;
    this._addedSignees = null;
    this._constitution = null;
    this._sendRules = null;

    this.dataKeyPrefix = "data_";
    this.removedSigneesKey = "removed_signees";
    this.addedSigneesKey = "added_signees";
    this.signeeUpdateLogKey = "signee_update_log";
    this.sendRulesKey = "send_rules";
    this.sendRulesUpdateLogKey = "send_rules_update_log";
    this.constitutionKey = "constitution";
    this.constitutionUpdateLogKey = "sys_config_update_log";

    this.actionDeleteSignee = "delete-signee";
    this.actionAddSignee = "add-signee";
    this.actionReplaceSignee = "replace-signee";
    this.actionSend = "send";
    this.actionUpdateConstitution = "update-constitution";
    this.actionUpdateSendRules = "update-rules";

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

    _getConstitution: function () {
        if (!this._constitution) {
            this._constitution = this.data.get(this.constitutionKey);
        }
        if (!this._constitution) {
            this._constitution = {
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
        return this._constitution;
    },

    _getSendRules: function () {
        if (!this._sendRules) {
            this._sendRules = this.data.get(this.sendRulesKey);
        }
        if (!this._sendRules) {
            // TODO: Confirm the default rules when going online
            this._sendRules = {
                "version": "0",
                "rules": [
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
        return this._sendRules;
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

    _getRemovedSignees: function () {
        if (!this._removedSignees) {
            this._removedSignees = this.data.get(this.removedSigneesKey);
        }
        if (!this._removedSignees) {
            this._removedSignees = [];
        }
        return this._removedSignees;
    },

    _getAddedSignees: function () {
        if (!this._addedSignees) {
            this._addedSignees = this.data.get(this.addedSigneesKey);
        }
        if (!this._addedSignees) {
            this._addedSignees = [];
        }
        return this._addedSignees;
    },

    _getSignees: function () {
        if (this._signees == null) {
            let removedManagers = this._getRemovedSignees();
            let addedManagers = this._getAddedSignees();

            this._signees = [];
            for (let i = 0; i < this.signees.length; ++i) {
                this._signees.push(this.signees[i]);
            }
            for (let i = 0; i < addedManagers.length; ++i) {
                this._signees.push(addedManagers[i]);
            }
            for (let i = 0; i < removedManagers.length; ++i) {
                let index = this._signees.indexOf(removedManagers[i]);
                this._signees.splice(index, 1);
            }
        }
        return this._signees;
    },

    _isRemoved: function (address) {
        return this._getRemovedSignees().indexOf(address) >= 0;
    },

    _verifyAddresses: function (address) {
        for (let i = 0; i < arguments.length; ++i) {
            if (0 === Blockchain.verifyAddress(arguments[i])) {
                throw (address + ' is not a valid nas address.');
            }
        }
    },

    _verifyDeleteAddress: function (address) {
        let managers = this._getSignees();
        if (managers.indexOf(address) < 0) {
            throw (address + ' could not be found.');
        }
    },

    _verifyAddAddress: function (address) {
        let managers = this._getSignees();
        if (managers.indexOf(address) >= 0) {
            throw (address + ' is already a manager.');
        }
        if (this._isRemoved(address)) {
            throw (address + ' has been removed.');
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

    _deleteSigneeAddress(address) {
        let removedManagers = this._getRemovedSignees();
        removedManagers.push(address);
        this.data.put(this.removedSigneesKey, removedManagers);
    },

    _addSigneeAddress(address) {
        let addedManagers = this._getAddedSignees();
        addedManagers.push(address);
        this.data.put(this.addedSigneesKey, addedManagers);
    },

    /**
     * { "action": "removed-signees", "detail": "n1xxxxx" }
     */
    _deleteSignee: function (data, signers) {
        let address = data.detail;
        this._verifyAddresses(address);
        this._verifyDeleteAddress(address);
        this._deleteSigneeAddress(address);
        data.signers = signers;
        this._addLog(data, this.signeeUpdateLogKey);
    },

    /**
     * { "action": "add-signees", "detail": "n1xxxxx" }
     */
    _addSignee: function (data, signers) {
        let address = data.detail;
        this._verifyAddresses(address);
        this._verifyAddAddress(address);
        this._addSigneeAddress(address);
        data.signers = signers;
        this._addLog(data, this.signeeUpdateLogKey);
    },

    /**
     * { "action": "replace-signees", "detail": { "oldAddress": "n1xxxxx", "newAddress": "n1yyyy" } }
     */
    _replaceSignee: function (data, signers) {
        let oldAddress = data.detail.oldAddress;
        let newAddress = data.detail.newAddress;
        this._verifyAddresses(oldAddress, newAddress);
        if (oldAddress === newAddress) {
            throw ('Old address cannot be the same as new address')
        }
        this._verifyDeleteAddress(oldAddress);
        this._verifyAddAddress(newAddress);
        this._deleteSigneeAddress(oldAddress);
        this._addSigneeAddress(newAddress);

        data.signers = signers;
        this._addLog(data, this.signeeUpdateLogKey);
    },

    /**
     * { "action": "send", "detail": { "id": "xxx", "to": "n1xxxxx", "value": "0.001" } }
     */
    _send: function (data, signers) {
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
     *     "action": "update-rules",
     *     "detail": {
     *          "version": "0",
     *          "rules": [
     *              { "startValue": "0", "endValue": "0.5", "proportionOfSigners": "0.3" },
     *              { "startValue": "0.5", "endValue": "1", "proportionOfSigners": "0.5" },
     *              { "startValue": "1", "endValue": this.infinity, "proportionOfSigners": "1" }
     *          ]
     *     ]
     * }
     */
    _updateSendRules: function (data, signers) {
        this._checkNumbers(data.detail.version);

        let ver = parseFloat(data.detail.version);
        if (ver <= parseFloat(this._getSendRules().version)) {
            throw ('Version error.');
        }
        let rules = data.detail.rules;
        if (rules.length <= 0) {
            throw ('Rules cannot be empty.');
        }
        let v = new BigNumber("0");
        for (let i = 0; i < rules.length; ++i) {
            let r = rules[i];
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
        this._sendRules = rules;
        this.data.put(this.sendRulesKey, rules);
        data.signers = signers;
        this._addLog(data, this.sendRulesUpdateLogKey);
    },

    /**
     * {
     *     "action": "update-constitution",
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
    _updateConstitution: function (data, signers) {
        this._checkNumbers(data.detail.version);
        let ver = parseFloat(data.detail.version);
        if (ver <= parseFloat(this._getConstitution().version)) {
            throw ('Version error.');
        }
        let ps = data.detail.proportionOfSigners;
        this._checkProportions(
            ps.updateSysConfig, ps.updateSendNasRule,
            ps.addManager, ps.deleteManager, ps.replaceManager
        );
        let config = data.detail;
        this._constitution = config;
        this.data.put(this.constitutionKey, config);
        data.signers = signers;
        this._addLog(data, this.constitutionUpdateLogKey);
    },

    _getNumOfNeedsSignersWithSendNasConfig: function (data) {
        let v = new BigNumber(data.detail.value);
        let rules = this._getSendRules().rules;
        let n = 1;
        for (let i = 0; i < rules.length; ++i) {
            let r = rules[i];
            let start = new BigNumber(r.startValue);
            let end = r.endValue === this.infinity ? null : new BigNumber(r.endValue);
            if (v.compareTo(start) >= 0 && (end == null || end.compareTo(v) > 0)) {
                n = parseFloat(r.proportionOfSigners);
                break;
            }
        }
        return n * this._getSignees();
    },

    _getNumOfNeedsSigners: function (data) {
        let p = this._getConstitution().proportionOfSigners;
        switch (data.action) {
            case this.actionAddSignee:
                return this._getSignees().length * parseFloat(p.addManager);

            case this.actionDeleteSignee:
                return (this._getSignees().length - 1) * parseFloat(p.deleteManager);

            case this.actionReplaceSignee:
                return (this._getSignees().length - 1) * parseFloat(p.replaceManager);

            case this.actionUpdateConstitution:
                return this._getSignees().length * parseFloat(p.updateSysConfig);

            case this.actionUpdateSendRules:
                return this._getSignees().length * parseFloat(p.updateSendNasRule);

            case this.actionSend:
                return this._getNumOfNeedsSignersWithSendNasConfig(data);

            default:
                throw ('Action ' + data.action + ' is not supported.');
        }
    },

    _verifySign: function (item) {
        let managers = this._getSignees();
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
            case this.actionAddSignee:
                this._addSignee(data, signers);
                break;
            case this.actionDeleteSignee:
                this._deleteSignee(data, signers);
                break;
            case this.actionReplaceSignee:
                this._replaceSignee(data, signers);
                break;
            case this.actionSend:
                this._send(data, signers);
                break;
            case this.actionUpdateSendRules:
                this._updateSendRules(data, signers);
                break;
            case this.actionUpdateConstitution:
                this._updateConstitution(data, signers);
                break;
            default:
                throw ('Action ' + data.action + ' is not supported.');
        }
    },

    getSignees: function () {
        return this._getSignees();
    },

    getConstitution: function () {
        return this._getConstitution()
    },

    getSendRules: function () {
        return this._getSendRules();
    },

    getSigneeUpdateLogs: function () {
        return this._getLogs(this.signeeUpdateLogKey);
    },

    getSendRulesUpdateLogs: function () {
        return this._getLogs(this.sendRulesUpdateLogKey);
    },

    getConstitutionUpdateLogs: function () {
        return this._getLogs(this.constitutionUpdateLogKey);
    },

    execute: function () {
        let array = arguments;
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

    query: function () {
        if (arguments.length === 0) {
            throw ('Arguments error. ');
        }
        return this._getData(arguments[0]);
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
