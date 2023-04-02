angular.module('sampleApp', ['ui.tree', 'ui.bootstrap']);

var SAMPLE_TYPE_MAPPING = {
    "brewery": {
        "enabled": true,
        "properties": {
            "name": {
                "fields":[{"name":"name","type":"text","analyzer":"en"}],
                "display_order": 2
            },
            "address": {
                "enabled": true,
                "properties": {
                    "city": {
                        "fields":[{"name":"city","type":"text"}]
                    }
                },
                "display_order": 1
            }
        },
        "display_order":"10"
    },
    "beer":{
        "enabled": true,
        "properties": {
            "name": {
                "fields":[{"name":"name","type":"text"}]
            }
        },
        "display_order":"5"
    }
};

var SAMPLE_INDEX_MAPPING = {
    "types": SAMPLE_TYPE_MAPPING,
    "default_mapping": {
        "enabled": true,
    },
    "type_field": "_type",
    "default_type": "_default",
    "default_analyzer": "standard",
    "default_datetime_parser": "dateTimeOptional",
    "default_field": "_all",
    "analysis": {
        "analyzers": {},
        "char_filters": {},
        "tokenizers": {},
        "token_filters": {},
        "token_maps": {
            "hw": {
                "type": "custom",
                "tokens": [
                    "hello",
                    "there"
                ]
            }
        }
    }
};

angular.module('sampleApp').
    controller('SampleCtrl', ['$scope', '$http', '$log', '$uibModal',
    function($scope, $http, $log, $uibModal) {
        var imc = initBleveIndexMappingController(
            $scope, $http, $log, $uibModal,
            SAMPLE_INDEX_MAPPING,
            {
                analyzerNames: null,
                dateTypeParserNames: null
            });

        $scope.done = function() {
            var im = imc.indexMapping();
            console.log(im);
            console.log(JSON.stringify(im));
        }
    }]);

angular.module('sampleApp').
    controller('BleveAnalyzerModalCtrl', BleveAnalyzerModalCtrl).
    controller('BleveCharFilterModalCtrl', BleveCharFilterModalCtrl).
    controller('BleveTokenizerModalCtrl', BleveTokenizerModalCtrl).
    controller('BleveTokenFilterModalCtrl', BleveTokenFilterModalCtrl).
    controller('BleveWordListModalCtrl', BleveWordListModalCtrl);
