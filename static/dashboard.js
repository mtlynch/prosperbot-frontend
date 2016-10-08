var model = {};
var dashboardApp = angular.module('dashboardApp', [])
  .directive('lineGraph', function($parse) {
    return {
      restrict: 'E',
      replace: false,
      scope: {data: '=chartData'},
      link: function (scope, element, attrs) {
        // Set the dimensions of the canvas / graph
        var margin = {top: 30, right: 20, bottom: 30, left: 50};
        var width = 900 - margin.left - margin.right;
        var height = 450 - margin.top - margin.bottom;

        // Set the ranges
        var x = d3.scaleTime().range([0, width]);
        var y = d3.scaleLinear().range([height, 0]);

        // Define the axes
        var xAxis = d3.axisBottom(x)
          .ticks(5);
        var yAxis = d3.axisLeft(y)
          .ticks(5);

        // Define the line
        var valueline = d3.line()
          .x(function(d) { return x(d.timestamp); })
          .y(function(d) { return y(d.value); });

        var updateGraph = function(data) {
          data.forEach(function(d) {
            d.timestamp = moment(d.Timestamp);
            d.value = scope.$eval(attrs['valueProperty'], d);
          });

          // Add the svg canvas
          var svg = d3.select(element[0])
            .append("svg")
              .attr("width", width + margin.left + margin.right)
              .attr("height", height + margin.top + margin.bottom)
          .append("g")
              .attr("transform",
                    "translate(" + margin.left + "," + margin.top + ")");

          // Scale the range of the data
          x.domain(d3.extent(data, function(d) { return d.timestamp; }));
          y.domain(d3.extent(data, function(d) { return d.value; }));

          // Add the X Axis
          svg.append("g")
            .attr("class", "x axis")
            .attr("transform", "translate(0," + height + ")")
            .call(xAxis);

          // Add the Y Axis
          svg.append("g")
           .attr("class", "y axis")
           .call(yAxis)
           .append("text")
           .attr("transform", "rotate(-90)")
           .attr("y", 6)
           .attr("dy", ".71em")
           .style("text-anchor", "end")
           .text("Balance ($)");

          // Add the valueline path.
          svg.append("path")
            .attr("class", "line")
            .attr("d", valueline(data));
        };

        scope.$watch('data', function(newValue) {
          if (!newValue) {
            return;
          }
          updateGraph(newValue);
        });
      }
    };
  });

dashboardApp.run(function($http) {
  var timeFormat = 'ddd, M/D/YY - h:mm:ss A';
  $http.get('/cashBalanceHistory').success(function(balances) {
    model.cashBalanceHistory = balances;
    model.latestCashBalance = balances[balances.length - 1];
  });
  $http.get('/accountValueHistory').success(function(totalValues) {
    model.accountValueHistory = totalValues;
    model.latestAccountValue = totalValues[totalValues.length - 1];
  });
  $http.get('/notes.json').success(function(notes) {
    model.notes = notes;
  });
});

dashboardApp.controller('DashboardCtrl', function ($scope) {
  $scope.dashboard = model;
});
