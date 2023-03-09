locals {
  cluster_id = var.eks-cluster.cluster_id
}
resource "datadog_dashboard_json" "environment_dashboard" {
  dashboard = <<EOF
  {
	"title": "[HAPPY] ${local.cluster_id} / Environment Dashboard",
	"description": "Created using the Datadog provider in Terraform",
	"widgets": [{
		"id": 3154357606055742,
		"definition": {
			"title": "Summary",
			"show_title": true,
			"type": "group",
			"layout_type": "ordered",
			"widgets": [{
				"id": 4805402064405576,
				"definition": {
					"title": "Response time  (avg)",
					"title_size": "13",
					"title_align": "left",
					"time": {
						"live_span": "1h"
					},
					"type": "query_value",
					"requests": [{
						"formulas": [{
							"formula": "query1 * 1000"
						}],
						"conditional_formats": [{
							"comparator": ">",
							"palette": "white_on_red",
							"value": 500
						}, {
							"comparator": ">",
							"palette": "white_on_yellow",
							"value": 400
						}, {
							"comparator": "<=",
							"palette": "white_on_green",
							"value": 400
						}],
						"response_format": "scalar",
						"queries": [{
							"query": "avg:aws.applicationelb.target_response_time.average{elbv2.k8s.aws/cluster:${local.cluster_id}}",
							"data_source": "metrics",
							"name": "query1",
							"aggregator": "avg"
						}]
					}],
					"autoscale": false,
					"custom_unit": "ms",
					"text_align": "left",
					"custom_links": [],
					"precision": 0
				},
				"layout": {
					"x": 0,
					"y": 0,
					"width": 2,
					"height": 2
				}
			}, {
				"id": 7100825526823894,
				"definition": {
					"title": "Healthy target count (min)",
					"title_size": "13",
					"title_align": "left",
					"time": {
						"live_span": "1h"
					},
					"type": "query_value",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"conditional_formats": [{
							"comparator": ">",
							"palette": "green_on_white",
							"value": 0
						}, {
							"comparator": "<=",
							"palette": "red_on_white",
							"value": 0
						}],
						"response_format": "scalar",
						"queries": [{
							"query": "sum:aws.applicationelb.healthy_host_count{elbv2.k8s.aws/cluster:${local.cluster_id}}",
							"data_source": "metrics",
							"name": "query1",
							"aggregator": "min"
						}]
					}],
					"autoscale": true,
					"custom_unit": "targets",
					"text_align": "left",
					"custom_links": [],
					"precision": 0
				},
				"layout": {
					"x": 2,
					"y": 0,
					"width": 2,
					"height": 2
				}
			}, {
				"id": 3655951002455680,
				"definition": {
					"title": "Unhealthy target count (max)",
					"title_size": "13",
					"title_align": "left",
					"time": {
						"live_span": "1h"
					},
					"type": "query_value",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"conditional_formats": [{
							"comparator": ">",
							"palette": "white_on_red",
							"value": 0
						}, {
							"comparator": "<=",
							"palette": "white_on_green",
							"value": 0
						}],
						"response_format": "scalar",
						"queries": [{
							"query": "sum:aws.applicationelb.un_healthy_host_count{elbv2.k8s.aws/cluster:${local.cluster_id}}",
							"data_source": "metrics",
							"name": "query1",
							"aggregator": "max"
						}]
					}],
					"autoscale": true,
					"custom_unit": "targets",
					"text_align": "left",
					"custom_links": [],
					"precision": 0
				},
				"layout": {
					"x": 4,
					"y": 0,
					"width": 2,
					"height": 2
				}
			}, {
				"id": 2809420673893146,
				"definition": {
					"title": "Requests per second (avg)",
					"title_size": "13",
					"title_align": "left",
					"time": {
						"live_span": "1h"
					},
					"type": "query_value",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "scalar",
						"queries": [{
							"query": "sum:aws.applicationelb.request_count{elbv2.k8s.aws/cluster:${local.cluster_id}}.as_rate()",
							"data_source": "metrics",
							"name": "query1",
							"aggregator": "avg"
						}]
					}],
					"autoscale": true,
					"text_align": "left",
					"custom_links": [],
					"precision": 1
				},
				"layout": {
					"x": 6,
					"y": 0,
					"width": 2,
					"height": 2
				}
			}]
		},
		"layout": {
			"x": 0,
			"y": 0,
			"width": 12,
			"height": 3
		}
	}, {
		"id": 2418827212694900,
		"definition": {
			"title": "Http Responses / Connections",
			"show_title": true,
			"type": "group",
			"layout_type": "ordered",
			"widgets": [{
				"id": 4642384245461586,
				"definition": {
					"title": "HTTP 2xx Responses",
					"title_size": "16",
					"title_align": "left",
					"show_legend": false,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {
						"live_span": "4h"
					},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "sum:aws.applicationelb.httpcode_target_2xx{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}.as_count()",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic"
						},
						"display_type": "bars"
					}],
					"custom_links": []
				},
				"layout": {
					"x": 0,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 5884051821562600,
				"definition": {
					"title": "HTTP 3xx Responses",
					"title_size": "16",
					"title_align": "left",
					"show_legend": false,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {
						"live_span": "4h"
					},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "sum:aws.applicationelb.httpcode_target_3xx{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}.as_count()",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic"
						},
						"display_type": "bars"
					}],
					"custom_links": []
				},
				"layout": {
					"x": 4,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 5108608634230402,
				"definition": {
					"title": "HTTP 4xx Responses",
					"title_size": "16",
					"title_align": "left",
					"show_legend": false,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {
						"live_span": "4h"
					},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"queries": [{
							"query": "sum:aws.applicationelb.httpcode_target_4xx{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}.as_count()",
							"data_source": "metrics",
							"name": "query1"
						}],
						"response_format": "timeseries",
						"style": {
							"palette": "warm"
						},
						"display_type": "bars"
					}],
					"custom_links": []
				},
				"layout": {
					"x": 8,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 2872253854507168,
				"definition": {
					"title": "HTTP 5xx Responses",
					"title_size": "16",
					"title_align": "left",
					"show_legend": false,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {
						"live_span": "4h"
					},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"queries": [{
							"query": "sum:aws.applicationelb.httpcode_target_5xx{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}.as_count()",
							"data_source": "metrics",
							"name": "query1"
						}],
						"response_format": "timeseries",
						"style": {
							"palette": "warm"
						},
						"display_type": "bars"
					}],
					"custom_links": []
				},
				"layout": {
					"x": 0,
					"y": 2,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 754659096536696,
				"definition": {
					"title": "Healthy Target Count",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "sum:aws.applicationelb.healthy_host_count{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 4,
					"y": 2,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 148794107811194,
				"definition": {
					"title": "Unhealthy Target Count",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "sum:aws.applicationelb.un_healthy_host_count{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 8,
					"y": 2,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 2466055870144974,
				"definition": {
					"title": "Active Connections",
					"title_size": "16",
					"title_align": "left",
					"show_legend": false,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {
						"live_span": "4h"
					},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "sum:aws.applicationelb.active_connection_count{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic"
						},
						"display_type": "bars"
					}],
					"custom_links": []
				},
				"layout": {
					"x": 0,
					"y": 4,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 8403694785490478,
				"definition": {
					"title": "New Connections",
					"title_size": "16",
					"title_align": "left",
					"show_legend": false,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {
						"live_span": "4h"
					},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "sum:aws.applicationelb.new_connection_count{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic"
						},
						"display_type": "bars"
					}],
					"custom_links": []
				},
				"layout": {
					"x": 4,
					"y": 4,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 7052214804631506,
				"definition": {
					"title": "Response Time",
					"title_size": "16",
					"title_align": "left",
					"show_legend": false,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {
						"live_span": "4h"
					},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1 * 1000"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "avg:aws.applicationelb.target_response_time.average{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "area"
					}],
					"custom_links": []
				},
				"layout": {
					"x": 0,
					"y": 6,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 61911063739598,
				"definition": {
					"title": "Processed Bytes",
					"title_size": "16",
					"title_align": "left",
					"show_legend": false,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"time": {
						"live_span": "4h"
					},
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "sum:aws.applicationelb.processed_bytes{elbv2.k8s.aws/cluster:${local.cluster_id}} by {happy_stack_name}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic"
						},
						"display_type": "area"
					}],
					"custom_links": []
				},
				"layout": {
					"x": 4,
					"y": 6,
					"width": 4,
					"height": 2
				}
			}]
		},
		"layout": {
			"x": 0,
			"y": 3,
			"width": 12,
			"height": 9
		}
	}, {
		"id": 2445078692824330,
		"definition": {
			"title": "Containers",
			"show_title": true,
			"type": "group",
			"layout_type": "ordered",
			"widgets": [{
				"id": 3679262507165410,
				"definition": {
					"title": "Container restarts",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "avg:kubernetes.containers.restarts{kube_cluster_name:${local.cluster_id},kube_namespace:${local.k8s_namespace}} by {happy_stack}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 0,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 7396801473522966,
				"definition": {
					"title": "Waiting containers",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "avg:kubernetes.containers.state.waiting{kube_cluster_name:${local.cluster_id},kube_namespace:${local.k8s_namespace}} by {happy_stack}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 4,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 5359542116641218,
				"definition": {
					"title": "Running containers in a stack",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "sum:kubernetes.containers.running{kube_cluster_name:${local.cluster_id},kube_namespace:${local.k8s_namespace}} by {happy_stack}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 8,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}]
		},
		"layout": {
			"x": 0,
			"y": 12,
			"width": 12,
			"height": 3
		}
	}, {
		"id": 3655088441997134,
		"definition": {
			"title": "Resources",
			"show_title": true,
			"type": "group",
			"layout_type": "ordered",
			"widgets": [{
				"id": 6368929839683538,
				"definition": {
					"title": "Stack Memory usage",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"alias": "Usage",
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "avg:kubernetes.memory.usage{kube_cluster_name:${local.cluster_id},kube_namespace:${local.k8s_namespace}} by {happy_stack}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 0,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 883128051754482,
				"definition": {
					"title": "Stack CPU Usage",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"alias": "Usage",
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "avg:kubernetes.cpu.usage.total{kube_cluster_name:${local.cluster_id},kube_namespace:${local.k8s_namespace}} by {happy_stack}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 4,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 1043953660409360,
				"definition": {
					"title": "Ephemeral Storage Usage",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"formula": "query1"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "avg:kubernetes.ephemeral_storage.usage{kube_cluster_name:${local.cluster_id},kube_namespace:${local.k8s_namespace}} by {happy_stack}",
							"data_source": "metrics",
							"name": "query1"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 8,
					"y": 0,
					"width": 4,
					"height": 2
				}
			}, {
				"id": 4906301887226174,
				"definition": {
					"title": "Network I/O",
					"title_size": "16",
					"title_align": "left",
					"show_legend": true,
					"legend_layout": "auto",
					"legend_columns": ["avg", "min", "max", "value", "sum"],
					"type": "timeseries",
					"requests": [{
						"formulas": [{
							"alias": "Read",
							"formula": "query1"
						}, {
							"alias": "Write",
							"formula": "query2"
						}],
						"response_format": "timeseries",
						"queries": [{
							"query": "avg:kubernetes.network.rx_bytes{kube_cluster_name:${local.cluster_id},kube_namespace:${local.k8s_namespace}} by {happy_stack}",
							"data_source": "metrics",
							"name": "query1"
						}, {
							"query": "avg:kubernetes.network.tx_bytes{kube_cluster_name:${local.cluster_id},kube_namespace:${local.k8s_namespace}} by {happy_stack}",
							"data_source": "metrics",
							"name": "query2"
						}],
						"style": {
							"palette": "dog_classic",
							"line_type": "solid",
							"line_width": "normal"
						},
						"display_type": "line"
					}]
				},
				"layout": {
					"x": 0,
					"y": 2,
					"width": 4,
					"height": 2
				}
			}]
		},
		"layout": {
			"x": 0,
			"y": 15,
			"width": 12,
			"height": 5
		}
	}],
	"template_variables": [],
	"layout_type": "ordered",
	"notify_list": [],
	"reflow_type": "fixed",
	"id": "ct9-mh6-fm8"
}
  EOF
}
