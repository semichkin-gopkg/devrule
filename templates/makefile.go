package templates

var Makefile = `
{{- define "BuildRule" -}}
{{- $command := .command | strings.TrimSpace | strings.ReplaceAll "\n" "\n\t" -}}
{{- if eq $command "" -}}
    {{ printf "%s: %s\n\n" .rule .dependencies }}
{{- else -}}
    {{ printf "%s: %s\n\t%s\n\n" .rule .dependencies $command }}
{{- end -}}
{{- end -}}

{{- define "ParseDict" -}}
{{- $result := dict -}}
{{- if has .dict .key -}}
    {{- $result = index .dict .key -}}
{{- end -}}
{{- $result | data.ToJSON -}}
{{- end -}}

{{- define "ParseSlice" -}}
{{- $result := coll.Slice -}}
{{- if has .dict .key -}}
    {{- $result = index .dict .key -}}
{{- end -}}
{{- $result | data.ToJSON -}}
{{- end -}}

{{- define "ParseString" -}}
{{- $result := "" -}}
{{- if has .dict .key -}}
    {{- $result = index .dict .key -}}
{{- end -}}
{{- $result | data.ToJSON -}}
{{- end -}}

{{/* Load configuration file (-d configuration.yaml) */}}
{{- $c := ds "configuration" -}}


{{/* Parse variables */}}
{{- $GV := tmpl.Exec "ParseDict" (dict "dict" $c "key" "GV") | data.JSON -}}
{{- $HR := tmpl.Exec "ParseDict" (dict "dict" $c "key" "HelperRules") | data.JSON -}}
{{- $ER := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "EnabledRules") | data.JSONArray -}}
{{- $DSR := tmpl.Exec "ParseDict" (dict "dict" $c "key" "DefaultServiceRules") | data.JSON -}}
{{- $S := tmpl.Exec "ParseDict" (dict "dict" $c "key" "Services") | data.JSON -}}


{{/* Render helper rules from configuration.HelperRules */}}
{{- "# HelperRules\n" -}}
{{- range $rule, $command := $HR -}}
    {{- template "BuildRule" dict "rule" $rule "dependencies" "" "command" $command }}
{{- end -}}


{{/* Render service rules from configuration.Services.[Name].Rules or configuration.DefaultServiceRules */}}
{{- "# ServiceRules\n" -}}
{{- range $name, $service := $S -}}

    {{/* Render service rules */}}
    {{- range $index, $rule := $ER -}}
        {{/* Load service rules */}}
        {{- $SR := tmpl.Exec "ParseDict" (dict "dict" $service "key" "Rules") | data.JSON -}}

        {{/* Generate command */}}
        {{- $command := "" -}}
        {{- if has $SR $rule -}}
            {{- $command = index $SR $rule -}}
        {{- else if has $DSR $rule -}}
            {{- $command = index $DSR $rule -}}
        {{- end -}}

        {{/* Load service variables */}}
        {{- $V := tmpl.Exec "ParseDict" (dict "dict" $service "key" "V") | data.JSON -}}

        {{/* Replace service variables inside command */}}
        {{- range $name, $value := $V -}}
            {{- $command = strings.ReplaceAll (printf "{{V.%s}}" $name) $value $command -}}
        {{- end -}}

        {{/* Replace global variables inside command */}}
        {{- range $name, $value := $GV -}}
            {{- $command = strings.ReplaceAll (printf "{{GV.%s}}" $name) $value $command -}}
        {{- end -}}

        {{- template "BuildRule" dict "rule" (printf "%s_%s" $name $rule) "dependencies" "" "command" $command -}}
    {{- end -}}

{{- end -}}

{{/* Render main rules */}}
{{- "# MainRules\n" -}}
{{- range $index, $rule := $ER -}}
    {{- $dependencies := "" -}}

    {{- range $name, $service := $S -}}
        {{- $dependencies = printf "%s %s_%s" $dependencies $name $rule -}}
    {{- end -}}

    {{- $dependencies = strings.TrimPrefix " " $dependencies}}
    {{- template "BuildRule" dict "rule" $rule "dependencies" $dependencies "command" "" -}}
{{- end -}}
`
