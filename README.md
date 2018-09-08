[![GoDoc](https://godoc.org/github.com/kuba--/splunk?status.svg)](http://godoc.org/github.com/kuba--/splunk)
[![Go Report Card](https://goreportcard.com/badge/github.com/kuba--/splunk)](https://goreportcard.com/report/github.com/kuba--/splunk)
[![Build Status](https://travis-ci.org/kuba--/splunk.svg?branch=master)](https://travis-ci.org/kuba--/splunk)
[![Version](https://badge.fury.io/gh/kuba--%2Fsplunk.svg)](https://github.com/kuba--/splunk/releases)
# splunk
The Splunk Enterprise REST API client

## Command line splunk search
The ```splunk``` search tool reads queries from *stdin* (line by line) and prints results on *stdout*.

### Usage
```
$ git clone git@github.com:kuba--/splunk.git
$ cd splunk
$ go install ./cmd/...
```

### Example
```
$ export SPLUNK_USERNAME=user
$ export SPLUNK_PASSWORD=pass

# Splunk API service runs on port :8089
$ export SPLUNK_URL=https://splunk.acme.com:8089

$ info

{
	"links": {},
	"origin": "https://splunk.acme.com:8089/services/server/info",
	"updated": "2016-11-08T17:13:48+00:00",
	"generator": {
		"build": "264376",
		"version": "6.2.3"
	},
	"entry": [
		{
			"name": "server-info",
			"id": "https://splunk.acme.com:8089/services/server/info/server-info",
			"updated": "2016-11-08T17:13:48+00:00",
			"links": {
				"alternate": "/services/server/info/server-info",
				"list": "/services/server/info/server-info"
			},
			"author": "system",
			"acl": {
				"app": "",
				"can_list": false,
				"can_write": false,
				"modifiable": false,
...


$ echo 'sourcetype="logs" host="*provisioning*" source="*.log"' | search -from -60min

# ... or multiple queries
$ echo 'sourcetype=src1 channel=service' > query
$ echo 'sourcetype=src2 host=*dev*' >> query
$ search -from -5min < query

{
	"preview": false,
	"offset": 9186,
	"result": {
		"_bkt": "service~6045~1AC69071-AC73-47C9-84E5-46AEDB65EACB",
		"_cd": "6045:261144704",
		"_indextime": "1477004266",
		"_raw": "{\"plug_idle\":false,\"ack_window\":0,\"duration\":844,\"thrift_process_start_ts\":1477004262587,\"consumer_src\":\"cm4S4TEOy9p3lrTYVU0MJ8L6KAh3AGbO\",\"plug_used_by\":\"derivative service backend\",\"start_time\":1477004262587,\"thrift_close_end_ts\":1477004263431,\"seq\":\"385555\",\"http_request_headers\":\"Accept-Encoding:gzip,User-Agent:Java/SDK/HttpClient,X-Forwarded-For:54.90.48.183\",\"plug_uptime_ms\":24886103,\"thrift_thread\":\"pool-16-thread-112\",\"plug_id\":\"jPnNlZ92A\",\"thrift_req_end_ts\":1477004263431,\"plug_checkpoint\":\"ZxDbrHb\",\"origin_server\":\"1cf10a76d2b0\",\"plug_ampq_broker\":0,\"api_level\":\"primary\",\"plug_type\":1,\"http_remoteip\":\"54.90.48.183\",\"api_category\":\"thrift/http\",\"http_method\":\"POST\",\"version\":\"1.0.0\",\"plug_container\":\"7e2396b02160\",\"source_service\":\"",\"api_scope\":\"F\",\"deployment\":\"teams.dev.pods.dev.us-east-8\",\"sdk_target\":\"teams.dev.pods.dev.us-east-8\",\"http_version\":\"1.0.0\",\"thrift_close_start_ts\":1477004263431,\"api_method\":\"/\",\"thrift_req\":\"pull_ack\",\"status\":\"ok\",\"events_bytes\":349,\"http_url\":\"eventing-dev.api.acme.com/\",\"plug_amqp_in\":10292,\"sdk_version\":\"java  SDK (1.0.16)\",\"play_thread_name\":\"play-akka.actor.default-dispatcher-13\",\"plug_build\":\"image-dfa7b7bd84\",\"facets_included\":\"http, ",\"plug_channel\":\"service\",\"thrift_code\":\"OK\",\"build_tag\":\"image-dfa7b7bd84\",\"out_seq\":\"385558\",\"plug_rollup\":36859693,\"finished_seq\":385559,\"payload_size\":150,\"thrift_process_end_ts\":1477004263431,\"plug_lag_ms\":820}",
		"_serial": "3701",
		"_si": [
			"splunk-dev-indexer-1",
			"service"
		],
		"_sourcetype": ",
		"_time": "2016-10-20 15:57:45.000 PDT",
		"plug_channel": [
			"service",
			"service"
		],
		"host": "1cf10a76d2b0",
		"index": "service",
		"linecount": "1",
		"source": "/var/log/stack-analytics.log",
		"sourcetype": "src1",
		"splunk_server": "splunk-dev-indexer-1"
	}
}
...
```
