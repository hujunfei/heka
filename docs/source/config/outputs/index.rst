.. _config_outputs:

=======
Outputs
=======

.. _config_common_output_parameters:

Common Output Parameters
========================

There are some configuration options that are universally available to all
Heka output plugins. These will be consumed by Heka itself when Heka
initializes the plugin and do not need to be handled by the plugin-specific
initialization code.

- message_matcher (string, optional):
    Boolean expression, when evaluated to true passes the message to the
    filter for processing. Defaults to matching nothing. See:
    :ref:`message_matcher`
- message_signer (string, optional):
    The name of the message signer. If specified only messages with this
    signer are passed to the filter for processing.
- ticker_interval (uint, optional):
    Frequency (in seconds) that a timer event will be sent to the filter.
    Defaults to not sending timer events.
- encoder (string, optional):
	Encoder to be used by the output. This should refer to the name of an
	encoder plugin section that is specified elsewhere in the TOML
	configuration. An instance of the specified encoder will be made available
	to the output plugin via the `Encoder()` method on the OutputRunner.
	Actual use of the specified encoder will vary from output to output.

.. _config_amqp_output:
.. include:: /config/outputs/amqp.rst

.. _config_carbon_output:
.. include:: /config/outputs/carbon.rst

.. _config_dashboard_output:
.. include:: /config/outputs/dashboard.rst

.. _config_elasticsearch_output:
.. include:: /config/outputs/elasticsearch.rst

.. _config_file_output:
.. include:: /config/outputs/file.rst

.. _config_log_output:
.. include:: /config/outputs/log.rst

.. _config_nagios_output:
.. include:: /config/outputs/nagios.rst

.. _config_smtp_output:
.. include:: /config/outputs/smtp.rst

.. _config_tcp_output:
.. include:: /config/outputs/tcp.rst

.. _config_whisper_output:
.. include:: /config/outputs/whisper.rst