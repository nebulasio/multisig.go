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

    this.sendDataKeyPrefix = "data_send_";
    this.voteDataKeyPrefix = "data_vote_";
    this.removedSigneesKey = "removed_signees";
    this.addedSigneesKey = "added_signees";
    this.signeeUpdateLogKey = "signee_update_log";
    this.sendRulesKey = "send_rules";
    this.sendRulesUpdateLogKey = "send_rules_update_log";
    this.constitutionKey = "constitution";
    this.constitutionUpdateLogKey = "constitution_update_log";

    this.actionDeleteSignee = "delete-signee";
    this.actionAddSignee = "add-signee";
    this.actionReplaceSignee = "replace-signee";
    this.actionUpdateConstitution = "update-constitution";
    this.actionUpdateSendRules = "update-rules";
    this.actionSend = "send";
    this.actionVote = "vote";

    this.voteAgree = "agree";
    this.voteDisagree = "disagree";
    this.voteValues = [this.voteAgree, this.voteDisagree, "abstain"];

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

    _setSendData: function (key, data) {
        this.data.put(this.sendDataKeyPrefix + key, data);
    },

    _getSendData: function (key) {
        return this.data.get(this.sendDataKeyPrefix + key);
    },

    _setVoteData: function (key, data) {
        this.data.put(this.voteDataKeyPrefix + key, data);
    },

    _getVoteData: function (key) {
        return this.data.get(this.voteDataKeyPrefix + key);
    },

    _getConstitution: function () {
        if (!this._constitution) {
            this._constitution = this.data.get(this.constitutionKey);
        }
        if (!this._constitution) {
            this._constitution = CONSTITUTION;
        }
        return this._constitution;
    },

    _getSendRules: function () {
        if (!this._sendRules) {
            this._sendRules = this.data.get(this.sendRulesKey);
        }
        if (!this._sendRules) {
            this._sendRules = RULES;
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
            if (n === undefined || n == null || !/^\d+(\.\d+)?$/.test(n)) {
                throw ('Data error.');
            }
        }
    },

    _checkProportions: function () {
        for (let i = 0; i < arguments.length; ++i) {
            let p = arguments[i];
            this._checkNumbers(p);
            p = parseFloat(p);
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
        if (this._getSendData(tx.id)) {
            throw (tx.id + ' exists.');
        }
        if (!Blockchain.transfer(tx.to, value)) {
            throw ('Transfer error.');
        }
        data.signers = signers;
        this._setSendData(tx.id, data);
    },

    /**
     * {
     *     "data": {
     *         "action": "vote",
     *         "detail": {
     *             "id": "xxxxx",
     *             "content": "test content",
     *             "proportionOfApproved": "0.8"
     *         }
     *     }
     *     "votes": {
     *         {
     *             "value": "agree", // abstain, agree, disagree
     *             "signer": "n1xxx",
     *             "sig": "1342abcdef...."
     *         }
     *         ...
     *     }
     * }
     */
    _vote: function (item) {
        let id = item.data.detail.id;
        if (!id) {
            throw ('vote data error. ');
        }

        if (this._getVoteData(id)) {
            throw ('The vote ' + id + ' exists.');
        }

        this._checkProportions(item.data.detail.proportionOfApproved);

        let votes = [];
        let n = 0;
        for (let i = 0; i < item.votes.length; ++i) {
            let vote = item.votes[i];
            if (this.voteValues.indexOf(vote.value) < 0) {
                throw ('Vote value "' + vote.value + '" error. ');
            }
            if (vote.value === this.voteAgree) {
                n++;
            }
            votes.push({"signer": vote.signer, "value": vote.value});
        }

        let t = parseFloat(item.data.detail.proportionOfApproved) * this._getSignees().length;
        let result = n >= t ? this.voteAgree : this.voteDisagree;
        if (n >= t) {
            let action = item.data.detail.approvedAction;
            if (action) {
                this._doVoteAction(action)
            }
        }
        this._setVoteData(id, {"data": item.data, "votes": votes, "result": result})
    },

    _doVoteAction: function (action) {
        switch (action.name) {
            case 'callContract':
                this._doCallContract(action.detail);
                break;
            default:
                throw ('Unknown action: ' + action.name);
        }
    },

    _doCallContract: function (detail) {
        this._verifyAddresses(detail.address);
        let args = detail.args;
        if (!detail.func || !args) {
            throw ('call contract "' + detail.address + '" function or args is null.');
        }
        args.splice(0, 0, detail.func);
        let c = new Blockchain.Contract(detail.address);
        c.value(0).call.apply(c, args);
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
     *              "updateConstitution": "1",
     *              "updateSendRules": "1",
     *              "addSignee": "1",
     *              "removeSignee": "1",
     *              "replaceSignee": "1",
     *              "vote": "1"
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
            ps.updateConstitution, ps.updateSendRules,
            ps.addSignee, ps.removeSignee, ps.replaceSignee,
            ps.vote
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
                return this._getSignees().length * parseFloat(p.addSignee);

            case this.actionDeleteSignee:
                return (this._getSignees().length - 1) * parseFloat(p.removeSignee);

            case this.actionReplaceSignee:
                return (this._getSignees().length - 1) * parseFloat(p.replaceSignee);

            case this.actionUpdateConstitution:
                return this._getSignees().length * parseFloat(p.updateConstitution);

            case this.actionUpdateSendRules:
                return this._getSignees().length * parseFloat(p.updateSendRules);

            case this.actionSend:
                return this._getNumOfNeedsSignersWithSendNasConfig(data);

            case this.actionVote:
                return this._getSignees().length * parseFloat(p.vote);

            default:
                throw ('Action ' + data.action + ' is not supported.');
        }
    },

    _verifySign: function (item) {
        let jsonData = JSON.parse(item.data);
        let action = jsonData.action;
        if ((action === this.actionVote && (item.votes === undefined || item.votes == null)) ||
            (action !== this.actionVote && (item.sigs === undefined || item.sigs == null))) {
            throw ('Data error. ')
        }

        let managers = this._getSignees();
        let signers = [];
        let n = this._getNumOfNeedsSigners(jsonData);

        if (action === this.actionVote) {
            if (item.votes.length < n) {
                throw ('Minimum ' + n + ' voters required.');
            }
            for (let j = 0; j < item.votes.length; ++j) {
                let v = item.votes[j];
                let hash = crypto.sha3256(item.data + v.value);
                let a = crypto.recoverAddress(1, hash, v.sig);
                if (a != null && managers.indexOf(a) >= 0) {
                    if (signers.indexOf(a) >= 0) {
                        throw ('Signature "' + v.sig + '" is repeated.');
                    } else {
                        signers.push(a);
                        v.signer = a;
                    }
                } else {
                    throw ('Signature "' + v.sig + '" error.');
                }
            }
        } else {
            if (item.sigs.length < n) {
                throw ('Minimum ' + n + ' signatures required.');
            }
            let hash = crypto.sha3256(item.data);
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
        }
        item.data = jsonData;
        return signers;
    },

    _execute: function (item, signers) {
        let data = item.data;
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
            case this.actionUpdateSendRules:
                this._updateSendRules(data, signers);
                break;
            case this.actionUpdateConstitution:
                this._updateConstitution(data, signers);
                break;
            case this.actionSend:
                this._send(data, signers);
                break;
            case this.actionVote:
                this._vote(item);
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
            this._execute(item, signers);
        }
        return "success";
    },

    querySentData: function () {
        if (arguments.length === 0) {
            throw ('Arguments error. ');
        }
        return this._getSendData(arguments[0]);
    },

    queryVotingData: function () {
        if (arguments.length === 0) {
            throw ('Arguments error. ');
        }
        return this._getVoteData(arguments[0]);
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
