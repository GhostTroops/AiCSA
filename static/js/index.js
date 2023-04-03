function IndexesCtrl($scope, $http, $routeParams, $log, $sce) {

	$scope.indexNames = [];
	$scope.errorMessage = null;

	$scope.clearErrorMessage = function() {
		$scope.errorMessage = null;
	};

	$scope.refreshIndexNames = function() {
		$http.get('/api').success(function(data) {
            $scope.indexNames = data.indexes;
        }).
        error(function(data, code) {
			$scope.errorMessage = data;
        });
	};

	$scope.deleteIndexNamed = function(name) {
		if(!confirm("Are you sure you want to permanenty delete the index '" + name + "'?")) {
			return;
		}
		$scope.clearErrorMessage();
		$http.delete('/api/' + name).success(function(data) {
			$scope.refreshIndexNames();
		}).
		error(function(data, code) {
			$scope.errorMessage = data;
		});
	};

	$scope.refreshIndexNames();
}

function IndexCtrl($scope, $http, $routeParams, $log, $sce, $uibModal) {

	$scope.indexName = $routeParams.indexName;
	$scope.indexDocCount = 0;
	$scope.indexDetails = null;
	$scope.mappingFormatted = "";
    $scope.viewOnly = true;

	$scope.tab = $routeParams.tabName;
	if($scope.tab === undefined || $scope.tab === "") {
		$scope.tab = "summary";
	}
	$scope.tabPath = '/static/partials/index/tab-' + $scope.tab + '.html';

	$scope.loadIndexDetails = function() {
		$http.get('/api/' + $scope.indexName).success(function(data) {
            $scope.indexDetails = data;
            $scope.mappingFormatted = JSON.stringify(data.mapping, undefined, 2);

            initBleveIndexMappingController($scope, $http, $log, $uibModal, data.mapping);
        }).
        error(function(data, code) {
			$scope.errorMessage = data;
        });
	};

	$scope.loadIndexDocCount = function() {
		$http.get('/api/' + $scope.indexName + '/_count').success(function(data) {
            $scope.indexDocCount = data.count;
        }).
        error(function(data, code) {
			$scope.errorMessage = data;
        });
	};

	// always load the details
	$scope.loadIndexDetails();

	// tab specific loading
	if($scope.tab === "summary") {
		$scope.loadIndexDocCount();
	}

	$scope.indexDocument = function(id, body) {
		$scope.clearErrorMessage();
		$http.put('/api/' + $scope.indexName + "/" + id, body).success(function(data) {
			$scope.successMessage = "Indexed Document: " + id;
		}).
		error(function(data, code) {
			$scope.errorMessage = data;
		});
	};

	$scope.deleteDocument = function(id) {
		$scope.clearErrorMessage();
		$http.delete('/api/' + $scope.indexName + "/" + id).success(function(data) {
			$scope.successMessage = "Deleted Document: " + id;
		}).
		error(function(data, code) {
			$scope.errorMessage = data;
		});
	};

    $scope.searchSyntax = function(query) {
        $http.post('/api/' + $scope.indexName + '/_search', {
            "size": 10,
            "explain": true,
            "highlight":{},
            "query": {
                "boost": 1.0,
                "query": query,
            }
        }).
        success(function(data) {
            $scope.processResults(data);
        }).
        error(function(data, code) {

        });
    };

    $scope.expl = function(explanation) {
            rv = "" + $scope.roundScore(explanation.value) + " - " + explanation.message;
            rv = rv + "<ul>";
            for(var i in explanation.children) {
                    child = explanation.children[i];
                    rv = rv + "<li>" + $scope.expl(child) + "</li>";
            }
            rv = rv + "</ul>";
            return rv;
    };

    $scope.roundScore = function(score) {
            return Math.round(score*1000)/1000;
    };

    $scope.roundTook = function(took) {
        if (took < 1000 * 1000) {
            return "less than 1ms";
        } else if (took < 1000 * 1000 * 1000) {
            return "" + Math.round(took / (1000*1000)) + "ms";
        } else {
            roundMs = Math.round(took / (1000*1000));
            return "" + roundMs/1000 + "s";
        }
	};

    $scope.processResults = function(data) {
        $scope.errorMessage = null;
        $scope.results = data;
        for(var i in $scope.results.hits) {
                hit = $scope.results.hits[i];
                hit.roundedScore = $scope.roundScore(hit.score);
                hit.explanationString = $scope.expl(hit.explanation);
                hit.explanationStringSafe = $sce.trustAsHtml(hit.explanationString);
                for(var ff in hit.fragments) {
                    fragments = hit.fragments[ff];
                    newFragments = [];
                    for(var ffi in fragments) {
                        fragment = fragments[ffi];
                        safeFragment = $sce.trustAsHtml(fragment);
                        newFragments.push(safeFragment);
                    }
                    hit.fragments[ff] = newFragments;
                }
        }
        $scope.results.roundTook = $scope.roundTook(data.took);
    };

    $scope.clearErrorMessage = function() {
		$scope.errorMessage = null;
	};
}

function makeIndexMapping() {
    return {
        "types": {},
        "default_mapping": {
            "enabled": true,
        },
        "type_field": "_type",
        "default_type": "_default",
        "default_analyzer": "standard",
        "default_datetime_parser": "dateTimeOptional",
        "default_field": "_all",
        "byte_array_converter": "json",
        "analysis": {
            "analyzers": {},
            "char_filters": {},
            "tokenizers": {},
            "token_filters": {},
            "token_maps": {}
        },
        "store_dynamic": true,
        "index_dynamic": true
    };
};

function IndexNewCtrl($scope, $http, $routeParams, $log, $location, $uibModal) {

	$scope.ierrorMessage = null;

    var imc = initBleveIndexMappingController(
        $scope, $http, $log, $uibModal, makeIndexMapping());

	$scope.newIndexNamed = function(name) {
		if (!name) {
			$scope.ierrorMessage = "Name is required";
			return;
		}

		$scope.clearErrorMessage();

		cleanMapping = imc.indexMapping();
		console.log(cleanMapping);

		$http.put('/api/' + name, cleanMapping).
            success(function(data) {
			    $location.path('/indexes/' + name);
		    }).
			error(function(data, code) {
				$scope.ierrorMessage = data;
			});
	};

	$scope.clearErrorMessage = function() {
		$scope.ierrorMessage = null;
	};
}
