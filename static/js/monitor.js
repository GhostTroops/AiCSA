function MonitorCtrl($scope, $http, $routeParams, $log, $sce, expvar) {

	$scope.active = true;
	$scope.monitoredIndexes = {};
	var updateInterval = null;

	var tv = 1000;

	var isoDateFormatter = function(x) { return ISODateString(new Date(x*1000)); };
	var byteSizeFormatter = function(y) { return Humanize.fileSize(y); };
	var nanosecondFormatter = function(y) {
		if(y > 1000000) {
			return Humanize.toFixed(y/1000000) + " ms";
		} else {
			return y + " ns";
		}
	};
	var intFormatter = function(y) {
		return Humanize.compactInteger(y);
	};

	var updateData = function() {
		expvar.pollExpvar();
		indxs = expvar.getKeys("indexes");
		indexesSeen = {};
		for(var idxIndex in indxs) {
			idxname = indxs[idxIndex];
			if ($scope.monitoredIndexes[idxname] === undefined) {
				// a new index
				$scope.monitoredIndexes[idxname] = monitorIndex(idxname);
			}
			indexesSeen[idxname] = true;
		}

		for (var monitoredIdxName in $scope.monitoredIndexes) {
			if (indexesSeen[monitoredIdxName] !== true) {
				$log.info("stop monitoring index: " + monitoredIdxName);
				removeIndex($scope.monitoredIndexes[monitoredIdxName]);
				delete $scope.monitoredIndexes[monitoredIdxName];
			}
		}

		for (var i in $scope.metrics) {
			category = $scope.metrics[i];
			for (var k in category.metrics) {
				metric = category.metrics[k];
				redrawMetric(metric);
			}
		}

		for (var idxName in $scope.monitoredIndexes) {
			idx = $scope.monitoredIndexes[idxName];
			if (idx !== undefined) {
				for(var j in idx.metrics) {
					metric = idx.metrics[j];
					redrawMetric(metric);
				}
			}
		}
	};

	function monitorIndex(name) {
		$log.info("start monitoring index: " + name);
		nameSelector = name.replace(".", "_")

		idx = {
			metrics: [
				{
					name: nameSelector+"updates",
					display: "Updates",
					path: "/bleve/indexes/" + name  + "/index/updates",
					type: "rate",
					color: "steelblue",
					yaxis: Rickshaw.Fixtures.Number.formatKMBT,
					xformatter: isoDateFormatter,
					yformatter: intFormatter
				},
				{
					name: nameSelector+"deletes",
					display: "Deletes",
					path: "/bleve/indexes/"+ name  +"/index/deletes",
					type: "rate",
					color: "steelblue",
					yaxis: Rickshaw.Fixtures.Number.formatKMBT,
					xformatter: isoDateFormatter,
					yformatter: intFormatter
				},
				{
					name: nameSelector+"indexanalysistime",
					display: "Analysis/Index Time",
					series: [
						{
							display: "Index",
							name: nameSelector+"indextime",
							path: "/bleve/indexes/"+ name  +"/index/index_time",
							type: "rate",
							color: "red",
						},
						{
							display: "Analysis",
							name: nameSelector+"analysistime",
							path: "/bleve/indexes/" + name  + "/index/analysis_time",
							type: "rate",
							color: "green",
						}
					],
					yaxis: Rickshaw.Fixtures.Number.formatKMBT,
					xformatter: isoDateFormatter,
					yformatter: nanosecondFormatter,
					legend: true
				},
				{
					name: nameSelector+"searches",
					display: "Searches",
					path: "/bleve/indexes/"+ name  +"/searches",
					type: "rate",
					color: "steelblue",
					yaxis: Rickshaw.Fixtures.Number.formatKMBT,
					xformatter: isoDateFormatter,
					yformatter: intFormatter
				},
				{
					name: nameSelector+"searchtime",
					display: "Search Time",
					path: "/bleve/indexes/"+ name  +"/search_time",
					type: "rate",
					color: "steelblue",
					yaxis: Rickshaw.Fixtures.Number.formatKMBT,
					xformatter: isoDateFormatter,
					yformatter: nanosecondFormatter
				},
		]
		};

		indexDivName = "index" + nameSelector;
		indexPanel = '<div id="panel'+nameSelector+'" class="panel panel-default"><div class="panel-heading"><a data-toggle="collapse" data-target="#' + indexDivName + '">' + name + '</a></div><div id="' + indexDivName + '" class="panel-body collapse in"></div></div>';
		$(indexPanel).insertBefore('#indexchartend');

		for (var i in idx.metrics) {
			metric = idx.metrics[i];

			divContent = '<h5 id="header' +metric.name+'">'+metric.display+'</h5><div id="'+metric.name+'"></div>';

			$(divContent).appendTo('#'+indexDivName);

			if (metric.legend) {
				legendDivName = "legend" + name;
				legendContent = '<div id="' + legendDivName + '" class="legend"></div>';
				$(legendContent).appendTo('#'+indexDivName);
			}

			// ask the expvar service to track this metric for us
			if (metric.path !== undefined) {
				expvar.addMetric(metric.name, metric.path);
			} else if (metric.series !== undefined) {
				for (var si in metric.series) {
					sm = metric.series[si];

					var swatch = document.createElement('div');
					swatch.className = 'swatch';
					swatch.style.backgroundColor = sm.color;
					$(swatch).appendTo('#'+legendDivName);

					var label = document.createElement('div');
					label.className = 'label';
					label.innerHTML = sm.display;
					$(label).appendTo('#'+legendDivName);

					expvar.addMetric(sm.name, sm.path);
				}
			}

			// build chart
			addGraph(metric);

		}

		return idx;
	}

	function removeIndex(index) {
		for(var i in index.metrics) {
			metric = index.metrics[i];
			expvar.removeMetric(metric.name);
			// $("#header" + metric.name).remove();
			// $("#"+metric.name).remove();
			$("#panel"+name).remove();
		}
	}

	function redrawMetric(metric) {
		graph = metric.graph;

		if (metric.series !== undefined) {
			var seriesData = [];
			for (var si in metric.series) {
				sm = metric.series[si];
				if (sm.type == "value") {
					currentValue = expvar.getMetricCurrentValue(sm.name);
					seriesData.push(currentValue);
				} else if (sm.type == "rate") {
					currentRate = expvar.getMetricCurrentRate(sm.name);
					seriesData.push(currentRate);
					$log.info("name: " + sm.name + " " + currentRate);
					$log.info(seriesData);
				}
			}
			graph.series.addData(seriesData);
		} else {
			var d = {};
			if (metric.type == "value") {
				currentValue = expvar.getMetricCurrentValue(metric.name);
				d[metric.name]= currentValue;
			} else if (metric.type == "rate") {
				currentRate = expvar.getMetricCurrentRate(metric.name);
				d[metric.name]= currentRate;
			}
			graph.series.addData(d);
		}
		

		// redraw
		graph.render();
	}

	// global metrics
	$scope.metrics = {
		"memory": {
			"display": "Memory",
			metrics: [
				{
					name: "alloc",
					display: "Memory Allocated",
					path: "/memstats/Alloc",
					type: "value",
					color: "steelblue",
					yaxis: Rickshaw.Fixtures.Number.formatKMBT,
					xformatter: isoDateFormatter,
					yformatter: byteSizeFormatter
				},
				{
					name: "pauseTotalNs",
					display: "Garbage Collection Time",
					path: "/memstats/PauseTotalNs",
					type: "rate",
					color: "steelblue",
					yaxis: Rickshaw.Fixtures.Number.formatKMBT,
					xformatter: isoDateFormatter,
					yformatter: nanosecondFormatter,
				}
			]
		}
	};


	expvar.addKeysLookup("indexes", "/bleve/indexes");

	for (var categoryName in $scope.metrics) {
		category = $scope.metrics[categoryName];

		divName = "cat" + categoryName;
		panel = '<div class="panel panel-default"><div class="panel-heading"><a data-toggle="collapse" data-target="#' + divName + '">'+ category.display + '</a></div><div id="' + divName + '" class="panel-body collapse in"></div></div>';
		$(panel).insertBefore('#chartend');


		for (var i in category.metrics) {
			metric = category.metrics[i];

			divContent = '<h5>'+metric.display+'</h5><div id="'+metric.name+'"></div>';

			$(divContent).appendTo("#"+divName);

			// ask the expvar service to track this metric for us
			expvar.addMetric(metric.name, metric.path);

			// build chart
			addGraph(metric);

		}
	}

	function addGraph(metric) {

		var seriesData = [];
		if (metric.series !== undefined) {
			for (var si in metric.series) {
				sm = metric.series[si];
				seriesData.push({
					name: sm.name,
					color: sm.color
				});
			}
		} else {
			seriesData.push({
				name: metric.name,
				color: metric.color
			});
		}

		$log.info("seriesdata");
		$log.info(seriesData);
		$log.info("seriesdataend");

		var graph = new Rickshaw.Graph({
			element: document.querySelector('#'+metric.name),
			width: "800",
			height: "75",
			renderer: "area",
			series: new Rickshaw.Series.FixedDuration(seriesData,
			undefined,
			{
				timeInterval: tv,
				maxDataPoints: 600,
				timeBase: new Date().getTime() / 1000
			})
		});

		// store the graph object inside the metric
		metric.graph = graph;

		// y-axis ticks
		if (metric.yaxis) {
			var yAxis = new Rickshaw.Graph.Axis.Y({
				graph: graph,
				tickFormat: metric.yaxis,
			});

			yAxis.render();
		}

		var xAxis = new Rickshaw.Graph.Axis.X({
			graph: graph,
			pixelsPerTick: 1000
		});
		xAxis.render();

		// set up the hover
		var hoverDetail = new Rickshaw.Graph.HoverDetail( {
			graph: graph,
			formatter: function(series, x, y, formattedX, formattedY, d) {
				var date = '<span class="x">' + formattedX + '</span>';
				var content =  formattedY + '<br>' + date;
				return content;
			}
		});

		if (metric.xformatter) {
			hoverDetail.xFormatter = metric.xformatter;
		}
		if (metric.yformatter) {
			hoverDetail.yFormatter = metric.yformatter;
		}

		// render it
		graph.render();
	}

	// setup data updates
	updateInterval = setInterval(updateData, tv);
	$scope.$on("$destroy", function(){
        clearInterval(updateInterval);
    });

	function ISODateString(d){
		function pad(n){return n<10 ? '0'+n : n;}
		return d.getUTCFullYear()+'-' +
			pad(d.getUTCMonth()+1)+'-' +
			pad(d.getUTCDate())+'T' +
			pad(d.getUTCHours())+':' +
			pad(d.getUTCMinutes())+':' +
			pad(d.getUTCSeconds())+'Z';
	}
}