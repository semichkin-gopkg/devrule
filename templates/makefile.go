package templates

const Makefile = `
{{- define "RenderRule" -}}
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

{{/* Load configuration file (-d configuration.yaml) */}}
{{- $c := ds "configuration" -}}

{{/* Init internal global rules */}}
{{- $IGR := dict "_clone" "[ -d '${to}' ] || git clone ${repo} ${to}" -}}


{{/* Parse variables */}}
{{- $GV := tmpl.Exec "ParseDict" (dict "dict" $c "key" "GlobalVars") | data.JSON -}}
{{- $GR := tmpl.Exec "ParseDict" (dict "dict" $c "key" "GlobalRules") | data.JSON -}}
{{- $MR := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "MainRules") | data.JSONArray -}}
{{- $DSR := tmpl.Exec "ParseDict" (dict "dict" $c "key" "DefaultServiceRules") | data.JSON -}}
{{- $S := tmpl.Exec "ParseDict" (dict "dict" $c "key" "Services") | data.JSON -}}


{{/* Render global rules */}}
{{- "# GlobalRules\n" -}}

{{/* Render internal global rules */}}
{{- range $rule, $command := $IGR -}}
    {{- template "RenderRule" dict "rule" $rule "dependencies" "" "command" $command }}
{{- end -}}

{{/* Render external global rules */}}
{{- range $rule, $command := $GR -}}
	{{/* Replace global variables inside command */}}
	{{- range $name, $value := $GV -}}
		{{- $command = strings.ReplaceAll (printf "{{GV.%s}}" $name) $value $command -}}
	{{- end -}}
    {{- template "RenderRule" dict "rule" $rule "dependencies" "" "command" $command }}
{{- end -}}


{{/* Render service rules from configuration.Services.[Name].Rules or configuration.DefaultServiceRules */}}
{{- "# ServiceRules\n" -}}
{{- range $name, $service := $S -}}

    {{/* Render service rules */}}
    {{- range $index, $rule := $MR -}}
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
        {{- $V = merge $V (dict "ServiceName" $name "Rule" $rule) -}}

        {{/* Replace service variables inside command */}}
        {{- range $name, $value := $V -}}
            {{- $command = strings.ReplaceAll (printf "{{V.%s}}" $name) $value $command -}}
        {{- end -}}

        {{/* Replace global variables inside command */}}
        {{- range $name, $value := $GV -}}
            {{- $command = strings.ReplaceAll (printf "{{GV.%s}}" $name) $value $command -}}
        {{- end -}}

        {{- template "RenderRule" dict "rule" (printf "%s_%s" $name $rule) "dependencies" "" "command" $command -}}
    {{- end -}}

{{- end -}}


{{/* Render main rules */}}
{{- "# MainRules\n" -}}
{{- range $index, $rule := $MR -}}
    {{- $dependencies := "" -}}

    {{- range $name, $service := $S -}}
        {{- $dependencies = printf "%s %s_%s" $dependencies $name $rule -}}
    {{- end -}}

    {{- $dependencies = strings.TrimPrefix " " $dependencies}}
    {{- template "RenderRule" dict "rule" $rule "dependencies" $dependencies "command" "" -}}
{{- end -}}
`
