#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
BUNDLE_DIR="${1}"
RBAC_DIR="${2}"
HELM_DIR="${3}"

if [[ "$OSTYPE" == "darwin"* ]]; then
  SEDPRG="gsed"
  REALPATHPRG="grealpath"
else
  SEDPRG="sed"
  REALPATHPRG="realpath"
fi

cd "${SCRIPT_DIR}"/../

# Split the generated bundle yaml file to inject control flags
yq e -Ns "\"${HELM_DIR}/templates/crds/\" + .spec.names.singular" ${BUNDLE_DIR}/bundle.yaml

# Add helm if statement for controlling the install of CRDs
for i in "${HELM_DIR}"/templates/crds/*.yml; do
  cp "$i" "$i.bkp"
  echo "# DO NOT EDIT, this file is autogenerated VIA helm.generate.sh!" > "$i"
  echo "{{- if .Values.installCRDs }}" >> "$i"
  cat "$i.bkp" >> "$i"
  rm "$i.bkp"
  $SEDPRG -i '0,/annotations/!b;//a\    {{- with .Values.crds.annotations }}\n    {{- toYaml . | nindent 4}}\n    {{- end }}\n' "$i"

  echo "{{- end }}" >> "$i"
  mv "$i" "${i%.yml}.yaml"
done

cat << EOF > ${HELM_DIR}/templates/rbac/kustomization.yaml
resources:
  - $($REALPATHPRG -s --relative-to ${HELM_DIR}/templates/rbac ${RBAC_DIR})
namePrefix: PLACEHOLDER-
EOF

kustomize build --load-restrictor LoadRestrictionsNone "${HELM_DIR}/templates/rbac" > ${HELM_DIR}/templates/rbac/rbac.yaml
rm ${HELM_DIR}/templates/rbac/kustomization.yaml

# Split the generated bundle yaml
yq e -Ns "\"${HELM_DIR}/templates/rbac/\" + .kind + \"-\" + (.metadata.name | sub(\"PLACEHOLDER-template-controller-\", \"\"))" ${HELM_DIR}/templates/rbac/rbac.yaml
rm ${HELM_DIR}/templates/rbac/rbac.yaml

# These are manually maintained
rm ${HELM_DIR}/templates/rbac/ServiceAccount-manager.yml
rm ${HELM_DIR}/templates/rbac/ClusterRoleBinding-manager-rolebinding.yml
rm ${HELM_DIR}/templates/rbac/RoleBinding-leader-election-rolebinding.yml

for i in "${HELM_DIR}"/templates/rbac/*.yml; do
  echo "# DO NOT EDIT, this file is autogenerated VIA helm.generate.sh!" > "$i.new"
  cat "$i" >> $i.new
  if [[ $i == *Role-*.yml ]]; then
    yq -i '.metadata.namespace="{{ .Release.Namespace }}"' $i.new
  fi
  $SEDPRG -i 's/PLACEHOLDER-template-controller/{{ include "template-controller.fullname" . }}/g' $i.new
  mv $i.new $i
done
