package templates

const Makefile = `
{{- define "RenderRuleWithoutNewLine" -}}
{{- $command := .command | strings.TrimSpace | strings.ReplaceAll "\n" "\n\t" -}}
{{- if eq $command "" -}}
    {{ printf "%s: %s\n" .rule .dependencies }}
{{- else -}}
    {{ printf "%s: %s\n\t%s\n" .rule .dependencies $command }}
{{- end -}}
{{- end -}}

{{- define "RenderRule" -}}
{{- template "RenderRuleWithoutNewLine" dict "rule" .rule "dependencies" .dependencies "command" .command -}}
{{ printf "\n" }}
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
{{- $EF := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "EnvFiles") | data.JSONArray -}}
{{- $GV := tmpl.Exec "ParseDict" (dict "dict" $c "key" "GV") | data.JSON -}}
{{- $GR := tmpl.Exec "ParseDict" (dict "dict" $c "key" "GlobalRules") | data.JSON -}}
{{- $MR := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "MainRules") | data.JSONArray -}}
{{- $DSR := tmpl.Exec "ParseDict" (dict "dict" $c "key" "DefaultServiceRules") | data.JSON -}}
{{- $S := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "Services") | data.JSONArray -}}


{{- /* Render env files import */ -}}
ifneq (,$(wildcard .env))
	include .env
	export
endif
ifneq (,$(wildcard .local.env))
	include .local.env
	export
endif
{{- range $i, $path := $EF }}
ifneq (,$(wildcard {{ $path }}))
	include {{ $path }}
	export
endif
{{- end }}

{{/* Render global rules */}}
{{- "\n# GlobalRules\n" -}}

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

{{- $Tags := coll.Slice "_" -}}

{{/* Render service rules from configuration.Services.[Name].Rules or configuration.DefaultServiceRules */}}
{{- "\n# ServiceRules\n" -}}
{{- range $index, $service := $S -}}

	{{- $serviceTags := tmpl.Exec "ParseSlice" (dict "dict" $service "key" "Tags") | data.JSONArray -}}
	{{- range $index, $tag := $serviceTags -}}
		{{- if not (has $Tags $tag) -}}
			{{- $Tags = $Tags | append $tag -}}
		{{- end -}}
	{{- end -}}

	{{- $name := index $service "Name" -}}
	{{/* Load service rules */}}
	{{- $SR := tmpl.Exec "ParseDict" (dict "dict" $service "key" "Rules") | data.JSON -}}

    {{/* Render service rules */}}
    {{- range $index, $rule := $MR -}}

        {{/* Generate command */}}
		{{- $commandName := printf "%s_%s" $name $rule -}}

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

		{{/* Render rule */}}
		{{- template "RenderRule" dict "rule" $commandName "dependencies" "" "command" $command -}}

    {{- end -}}

{{- end -}}

{{/* Render grouped rules */}}
{{- "\n# GroupedRules\n\n" -}}
{{- range $index, $tag := $Tags -}}
	{{- if eq $tag "_" -}}
		{{- "# Main Rules\n" -}}
	{{- else -}}
		{{- printf "# %s Rules\n" $tag -}}
	{{- end -}}

	{{- range $index, $rule := $MR -}}
		{{- $dependencies := "" -}}
	
		{{- range $index, $service := $S -}}
			{{- $serviceTags := tmpl.Exec "ParseSlice" (dict "dict" $service "key" "Tags") | data.JSONArray -}}

			{{- if or (eq $tag "_") (has $serviceTags $tag) -}}
				{{- $dependencies = printf "%s %s_%s" $dependencies (index $service "Name") $rule -}}
			{{- end -}}
		{{- end -}}

		{{- $dependencies = strings.TrimPrefix " " $dependencies}}

		{{- $groupRule := $rule -}}
		{{- if ne $tag "_" -}}
			{{- $groupRule = printf "%s_%s" $tag $rule -}}
		{{- end -}}

		{{- template "RenderRuleWithoutNewLine" dict "rule" $groupRule "dependencies" $dependencies "command" "" -}}
	{{- end -}}

	{{ printf "\n" }}
{{- end -}}
`
