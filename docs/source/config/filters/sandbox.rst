
SandboxFilter
=============

The sandbox filter provides an isolated execution environment for data
analysis. Any output generated by the sandbox is injected into the payload of
a new message for further processing or to be output.

Config:

- :ref:`config_common_filter_parameters`

- script_type (string):
    The language the sandbox is written in.  Currently the only valid option
    is 'lua'.

- filename (string):
    For a static configuration this is the path to the sandbox code; if
    specified as a relative path it will be appended to Heka's global
    share_dir. The filename must be unique between static plugins, since the
    global data is preserved using this name. For a dynamic configuration the
    filename is ignored and the the physical location on disk is controlled by
    the SandboxManagerFilter.

- preserve_data (bool):
    True if the sandbox global data should be preserved/restored on Heka
    shutdown/startup.

- memory_limit (uint):
    The number of bytes the sandbox is allowed to consume before being
    terminated (max 8MiB, default max). For a dynamic configuration the
    value is ignored and the SandboxManagerFilter setting is used.

- instruction_limit (uint):
    The number of instructions the sandbox is allowed the execute during the
    process_message/timer_event functions before being terminated (max 1M,
    default max). For a dynamic configuration the value is ignored and the 
    SandboxManagerFilter setting is used.

- output_limit (uint):
    The number of bytes the sandbox output buffer can hold before before being
    terminated (max 63KiB, default max).  Anything less than 64B is set to
    64B. For a dynamic configuration the value is ignored and the 
    SandboxManagerFilter setting is used.

- profile (bool):
    When true a statistically significant number of ProcessMessage timings are
    immediately captured before reverting back to the regular sampling
    interval.  The main purpose is for more accurate sandbox
    comparison/tuning/optimization.

- module_directory (string):
    The directory where 'require' will attempt to load the external Lua
    modules from.  Defaults to ${SHARE_DIR}/lua_modules. For a dynamic
    configuration the module_directory is ignored and the the physical
    location on disk is controlled by the SandboxManagerFilter.

- config (object):
    A map of configuration variables available to the sandbox via read_config.
    The map consists of a string key with: string, bool, int64, or float64
    values.

Example:

.. code-block:: ini

    [hekabench_counter]
    type = "SandboxFilter"
    message_matcher = "Type == 'hekabench'"
    ticker_interval = 1
    script_type  = "lua"
    filename = "counter.lua"
    preserve_data = true
    profile = false

        [hekabench_counter.config]
        rows = 1440
        sec_per_row = 60
