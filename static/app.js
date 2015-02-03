(function () {
	var app = angular.module('spinClass', ['ui.bootstrap']);

	app.filter('doubledstring', function() {
		return function(input) {
			string = String(input);
			if ( string.length == 1 ) {
				string = "0" + string
			}
			return string
		}
	});

	var SpinService = function($http, $interval) {
		this.$http = $http;
		this.$interval = $interval;
	}

	SpinService.prototype.spinUp = function(count) {
		var _this = this;
		return _this.$http.post('/spin/up', {count: count});
	}

	SpinService.prototype.spinDown = function(prefix) {
		var _this = this;
		return _this.$http.post('/spin/down', {prefix: prefix})
	}

	SpinService.prototype.final = function(prefix) {
		var _this = this;
		return _this.$http.post('/final', {prefix: prefix})
	}

	SpinService.prototype.odometer = function(prefix) {
		var _this = this;
		return _this.$http.post('/odometer', {prefix: prefix});
	}

	SpinService.$inject = ['$http', '$interval']
	app.service('spinService', SpinService)


	var Spin = function(spinService, $interval) {
		this.$interval = $interval
		this.spinService = spinService
		this.total = this.building = this.success = this.failed = this.minutes = this.seconds = 0;

		this.statusMap = {
			true: {
				'ACTIVE': 'warning',
				'DELETED': 'success',
				'': 'success',
			},
			false: {
				'BUILD': 'warning',
				'ACTIVE': 'success',
				'ERROR': 'danger',
				'': 'info'
			}
		}
	}

	Spin.prototype.action = function() {
		var _this = this;
		if (_this.down) {
			return 'Deleting'
		}
		return 'Building'
	}

	Spin.prototype.showStats = function() {
		var _this = this;
		if (_this.down === true) {
			return false;
		}
		if (_this.prefix != "") {
			return true;
		}
	}

	Spin.prototype.percent = function(value) {
		var _this = this;
		if ( !_this.total || !value ) {
			return 0;
		}
		return Math.round((parseInt(value) / _this.total) * 100);
	}

	Spin.prototype.odometer = function(prefix, down) {
		var _this = this;

		if ( down ) {
			_this.seconds--;
			if ( _this.seconds == 0 ) {
				_this.minutes--;
				_this.seconds = 59;
			}

		} else {
			_this.seconds++;
			if ( _this.seconds == 60 ) {
				_this.minutes++;
				_this.seconds = 0;
			}
		}

		_this.spinService.odometer(prefix).
			success(function(data) {
				_this.failed = 0;
				_this.success = 0;
				_this.building = 0;
				_this.instances = data.instances;
				for ( id in data.instances ) {
					var instance = data.instances[id]
					data.instances[id].StatusType = _this.statusMap[down][instance.Status];
					if ( instance.Status == "ERROR" ) {
						_this.failed++;
					} else if (instance.Status == "ACTIVE" && !down) {
						_this.success++;
					} else if ((instance.Status == 'DELETED' || instance.Status == '') && down) {
						_this.success++;
					} else {
						_this.building++;
					}
				}
				if ( _this.failed + _this.success == _this.total ) {
					_this.$interval.cancel(_this.odometerPromise);
					_this.odometerPromise = false;
					if ( down ) {
						_this.spinService.final(prefix);
						_this.prefix = '';
						_this.instances = {};
						_this.down = false;
						_this.total = _this.success = _this.failed = _this.building = _this.minutes = _this.seconds = 0;
					}
				}
			});
	}

	Spin.prototype.spinUp = function(count) {
		var _this = this;

		count = parseInt(count)
		_this.total = count;
		if ( count == 0 ) {
			return
		}

		_this.spinService.spinUp(count).
			success(function(data) {
				_this.prefix = data.prefix;
				_this.odometerPromise = _this.$interval(function() { _this.odometer(data.prefix, false) }, 1000)
				_this.counterPromise = _this.$interval(_this.countUp, 1000)
			});
	}

	Spin.prototype.spinDown = function(prefix) {
		var _this = this;

		if ( _this.odometerPromise != false ) {
			_this.$interval.cancel(_this.odometerPromise);
			_this.odometerPromise = false;
		}
		_this.down = true;
		_this.spinService.spinDown(prefix).
			success(function(data) {
				_this.odometerPromise = _this.$interval(function() { _this.odometer(prefix, true) }, 1000)
			});
	}

	Spin.$inject = ['spinService', '$interval']

	app.controller('Spin', Spin);
})()
