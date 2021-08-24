![GitHub](https://img.shields.io/github/license/leroy-merlin-br/logstash-exporter) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/leroy-merlin-br/logstash-exporter)

# logstash-exporter 

Exports [Logstash]() metrics to [Prometheus]() for monitoring.

## Version compatibility and some important info

The exporter only supports Logstash >= v7.3. For earlier versions please check 
[sequra/logstash_exporter](https://github.com/sequra/logstash_exporter).

As some may notice, this is exporter is a result of multiple forks, we decided to create another version for a number 
of reasons:
- [sequra/logstash_exporter](https://github.com/sequra/logstash_exporter) maintainers doesn't seem to care about the 
  project anymore.
- This project builds binaries and docker images, eliminating the need for a project like
  [bitnami/bitnami-docker-logstash-exporter](https://github.com/bitnami/bitnami-docker-logstash-exporter) which downloads
  a binary from [their servers](https://github.com/bitnami/bitnami-docker-logstash-exporter/blob/master/7.3/debian-10/Dockerfile#L12), 
  which is **super creepy**.
- Bitnami's project doesn't offer a build for other architectures such as `arm64` and here at Leroy Merlin Brasil we 
  rely heavily on Logstash pods running on Graviton nodes.

## Run

For pre-built binaries, check the [releases page](https://github.com/leroy-merlin-br/logstash-exporter/releases).

### Docker

```sh
docker run leroymerlinbr/logstash_exporter:latest
```

### CLI

```sh
logstash_exporter <flags>
  -h, --help                    Show context-sensitive help (also try --help-long and --help-man).
      --logstash.endpoint="http://localhost:9600"  
                                The protocol, host and port on which logstash metrics API listens.
      --listen.address=":9198"  Address on which to expose metrics and web interface.
      --log.level="info"        The logging level to be defined.
      --version                 Show application version.
```

## Metrics

| Name | Type | Description |
| --- | --- | --- |
| `logstash_exporter_build_info` | gauge | Exporter build info |
| `logstash_exporter_scrape_duration_seconds` | summary | Duration of a scrape job. |
| `logstash_info_jvm` | counter |  A metric with a constant '1' value labeled by name, version and vendor of the JVM running Logstash.| 
| `logstash_info_node`| counter |  A metric with a constant '1' value labeled by Logstash version. |
| `logstash_info_os` | counter | A metric with a constant '1' value labeled by name, arch, version and available_processors to the OS running Logstash. |
| `logstash_node_gc_collection_duration_seconds_total` | counter | |
| `logstash_node_gc_collection_total` | gauge | | 
| `logstash_node_jvm_threads_count` | gauge | |
| `logstash_node_jvm_threads_peak_count` | gauge | | 
| `logstash_node_mem_heap_committed_bytes` | gauge | |
| `logstash_node_mem_heap_max_bytes` | gauge | |
| `logstash_node_mem_heap_used_bytes` | gauge | 
| `logstash_node_mem_nonheap_committed_bytes` | gauge | |
| `logstash_node_mem_nonheap_used_bytes` | gauge | |
| `logstash_node_mem_pool_committed_bytes` | gauge | | 
| `logstash_node_mem_pool_max_bytes` | gauge | |
| `logstash_node_mem_pool_peak_max_bytes` | gauge | |
| `logstash_node_mem_pool_peak_used_bytes` | gauge | |
| `logstash_node_mem_pool_used_bytes` | gauge | |
| `logstash_node_pipeline_duration_seconds_total` | counter | |
| `logstash_node_pipeline_events_filtered_total` | counter | |
| `logstash_node_pipeline_events_in_total` | counter | |
| `logstash_node_pipeline_events_out_total` | counter | |
| `logstash_node_pipeline_queue_push_duration_seconds_total` | counter | |
| `logstash_node_plugin_bulk_requests_failures_total` | counter | |
| `logstash_node_plugin_bulk_requests_successes_total` | counter | |
| `logstash_node_plugin_bulk_requests_with_errors_total` | counter | |
| `logstash_node_plugin_documents_failures_total` | counter | |
| `logstash_node_plugin_documents_successes_total` | counter | |
| `logstash_node_plugin_duration_seconds_total` (| counter | |
| `logstash_node_plugin_queue_push_duration_seconds_total` | counter | |
| `logstash_node_plugin_events_in_total` | counter | |
| `logstash_node_plugin_events_out_total` | counter | |
| `logstash_node_plugin_current_connections_count` | gauge
| `logstash_node_plugin_peak_connections_count` | gauge
| `logstash_node_process_cpu_total_seconds_total` | counter | |
| `logstash_node_process_max_filedescriptors` | gauge
| `logstash_node_process_mem_total_virtual_bytes` | gauge
| `logstash_node_process_open_filedescriptors` | gauge
| `logstash_node_queue_events` | counter | |
| `logstash_node_queue_size_bytes` | counter | |
| `logstash_node_queue_max_size_bytes` | counter | |
| `logstash_node_dead_letter_queue_size_bytes` | counter | |
| `logstash_node_up`: | gauge | whether logstash node is up (1) or not (0) |

## Contributing

We welcome any contributions. We appreciate pretty git commit messages such as 
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

Please check the [Issues](https://github.com/leroy-merlin-br/logstash-exporter/issues) page for features and bugs that
need some love.

### Integration tests

In order to execute manual integration tests (to be sure certain logstash version is compatible with logstash_exporter),
you can follow instructions [here](integration-tests/README.md).
