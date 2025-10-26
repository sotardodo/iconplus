#!/bin/bash

mkdir -p /var/lib/rundeck/jobs/{{ rundeck.project_name }}-{{ envDeployment }}/{{ item.name }}/config

{% for volume in item.volumes %}
mkdir -p /var/lib/rundeck/jobs/{{ rundeck.project_name }}-{{ envDeployment }}/{{ item.name }}/volumes/{{ volume.name }}
{% endfor %}

{% for config in configFiles %}
echo "Note from DevOps Team, at (@) character will be escape with double at. Rundeck have limitation on that character :)"
echo {{ lookup('file', tempLocation + config.dest) | replace("@", "@@") | quote }} > /var/lib/rundeck/jobs/{{ rundeck.project_name }}-{{ envDeployment }}/{{ item.name }}/config/{{ config.dest }}
sed -i -e 's/@@/@/g' /var/lib/rundeck/jobs/{{ rundeck.project_name }}-{{ envDeployment }}/{{ item.name }}/config/{{ config.dest }}
{% endfor %}

