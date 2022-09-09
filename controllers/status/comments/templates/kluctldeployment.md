
{% set target = get_var("object.status.target", "") | from_yaml %}
{% set conditionsByType = dict(get_var("object.status.conditions", []) | groupby(attribute="type")) %}

{% macro printErrors(l) %}
| kind | namespace/name | message |
|------|----------------|---------|
{% for e in l %}
| {{ e.ref.kind }} | {{ e.ref.namespace or "<global>" }}/{{ e.ref.name }} | {{ e.error if e.error is defined else e.message }} |
{% endfor %}
{% endmacro %}

{% macro printResult(title, result, success_message, warning_message, error_message, logs_message, cmd_name, show_changes) %}
{% set decoded=result.rawResult|from_yaml if result.rawResult is defined else none %}
### {% if (result.error is defined and result.error != "") or (decoded.errors is defined and decoded.errors|length > 0) %}:boom:{% elif decoded.warnings is defined and decoded.warnings|length > 0 %}:warning:{% else %}:white_check_mark:{% endif %} {{ title }}
<details>
<summary>Click to expand</summary>

{% if (result.error is defined and result.error != "") or (decoded.errors is defined and decoded.errors|length > 0) %}
{{ error_message.format(time=result.time) }}
{% elif decoded.warnings is defined and decoded.warnings|length > 0 %}
{{ warning_message.format(time=result.time) }}
{% else %}
{{ success_message.format(time=result.time) }}
{% endif %}

{% if result.error is defined and result.error != "" %}
:boom: :boom: :boom: Command resulted in error: {{ result.error }}
{% endif %}

{{ logs_message.format(url="TODO") }}
{% if decoded %}
{% if show_changes %}
#### Changes
{% if decoded.get("newObjects") %}:floppy_disk: New K8s objects: {{ decoded.newObjects|length }}<br>{% endif %}
{% if decoded.get("changedObjects") %}:construction_worker: Changed K8s objects: {{ decoded.changedObjects|length }}<br>{% endif %}
{% if decoded.get("deletedObjects") %}:broken_heart: Deleted K8s objects: {{ decoded.deletedObjects|length }}<br>{% endif %}
{% endif %}
{% if decoded.errors is defined and decoded.errors|length != 0 %}
#### :boom: Errors
{{ printErrors(decoded.errors) }}
{% endif %}
{% if decoded.warnings is defined and decoded.warnings|length != 0 %}
#### :warning: Warnings
{{ printErrors(decoded.warnings) }}
{% endif %}
{% endif %}
</details>

{% if decoded.results is defined and decoded.results|length > 0 %}
# :tada: Results
| Key | Message |
|--------|---------|
{% for e in decoded.results %}
| {{ e.annotation.replace("validate-result.kluctl.io/", "") }} | {{ e.message }} |
{% endfor %}
{% endif %}

{% endmacro %}

# :robot: Deployment Summary ({{ object.spec.target }})

{% set result=get_var("object.status.lastDeployResult", none) %}
{% if result %}
{{ printResult("Deployment", result,
        "The deployment succeeded at {time}",
        "Warning(s) occurred at {time} while performing the deployment.",
        "Error(s) occurred at {time} while performing the deployment.",
        "Logs of the last deployment can be found [here]({url}).",
        "deploy",
        True) }}
{% endif %}

{% set result=get_var("object.status.lastValidateResult", none) %}
{% if result %}
{{ printResult("Validation", result,
        "The most recent validation succeeded",
        "Warning(s) occurred while validating the current deployment state.",
        "Error(s) occurred while validating the current deployment state.",
        "Logs of the last validation can be found [here]({url}).",
        "validate",
        False) }}
{% endif %}

{% if "Ready" in conditionsByType and conditionsByType["Ready"] %}
{% if conditionsByType["Ready"][0].status == "True" %}
# :robot: Target is ready! :thumbsup:
{% else %}
# :robot: Target is NOT ready! :thumbsdown:<br><br>
Message from controller: {{ conditionsByType["Ready"][0].message }}
{% endif %}
{% endif %}
