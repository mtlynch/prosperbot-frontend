var model = {};
var notesApp = angular.module('notesApp', []);

notesApp.run(function($http) {
  $http.get('/notes.json').success(function(notes) {
    model.notes = notes;
  });
});

notesApp.constant('NOTE_STATUS', {
  'ORIGINATION_DELAYED': 0,
  'CURRENT': 1,
  'CHARGEOFF': 2,
  'DEFAULTED': 3,
  'COMPLETED': 4,
  'FINAL_PAYMENT_IN_PROGRESS': 5,
  'CANCELLED': 6,
});

notesApp.filter('statusName', function(NOTE_STATUS) {
  return function (status) {
    names = {}
    names[NOTE_STATUS.CURRENT] = 'Current';
    names[NOTE_STATUS.CHARGEOFF] = 'Charged off'
    names[NOTE_STATUS.DEFAULTED] = 'Defaulted'
    names[NOTE_STATUS.COMPLETED] = 'Completed'
    return names[status];
  };
});

notesApp.controller('NotesCtrl', function ($scope, NOTE_STATUS) {
  $scope.model = model;
  $scope.getStatusClass = function (status) {
    classSuffix = {}
    classSuffix[NOTE_STATUS.CURRENT] = 'info';
    classSuffix[NOTE_STATUS.CHARGEOFF] = 'danger'
    classSuffix[NOTE_STATUS.DEFAULTED] = 'warning'
    classSuffix[NOTE_STATUS.COMPLETED] = 'success'
    return 'alert alert-' + classSuffix[status];
  }
});
