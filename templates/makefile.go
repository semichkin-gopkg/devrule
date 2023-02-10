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

{{/* Init internal rules */}}
{{- $IR := dict "_clone" "@[ -d '${to}' ] || git clone ${repo} ${to}" -}}


{{/* Parse variables */}}
{{- $EX := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "Expressions") | data.JSONArray -}}
{{- $EX = $EX | prepend "export PWD := $(shell pwd)" -}}

{{- $EF := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "EnvFiles") | data.JSONArray -}}
{{- $EF = $EF | prepend ".env" -}}

{{- $GV := tmpl.Exec "ParseDict" (dict "dict" $c "key" "GV") | data.JSON -}}
{{- $GR := tmpl.Exec "ParseDict" (dict "dict" $c "key" "GlobalRules") | data.JSON -}}
{{- $MR := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "MainRules") | data.JSONArray -}}
{{- $DSR := tmpl.Exec "ParseDict" (dict "dict" $c "key" "DefaultServiceRules") | data.JSON -}}
{{- $S := tmpl.Exec "ParseSlice" (dict "dict" $c "key" "Services") | data.JSONArray -}}

{{- /* Render expressions */ -}}
{{- "# Expressions\n" -}}
{{- range $i, $expression := $EX }}
{{- printf "%s\n" $expression -}}
{{- end -}}
{{- "\n\n" -}}

{{- /* Render env files import */ -}}
{{- "# EnvFiles\n" -}}
{{- range $i, $path := $EF -}}
ifneq (,$(wildcard {{ $path }}))
	include {{ $path }}
	export
endif
{{ end }}

{{/* Render internal rules */}}
{{- "# InternalRules\n" -}}
{{- range $rule, $command := $IR -}}
    {{- template "RenderRule" dict "rule" $rule "dependencies" "" "command" $command }}
{{- end -}}
{{- "\n" -}}

{{/* Render global rules */}}
{{- "# GlobalRules\n" -}}
{{- range $rule, $command := $GR -}}
	{{/* Replace global variables inside command */}}
	{{- range $name, $value := $GV -}}
		{{- $command = strings.ReplaceAll (printf "{{GV.%s}}" $name) $value $command -}}
	{{- end -}}
    {{- template "RenderRule" dict "rule" $rule "dependencies" "" "command" $command }}
{{- end -}}
{{- "\n" -}}

{{- $Groups := coll.Slice "_" -}}

{{/* Render service rules from configuration.Services.[Name].Rules or configuration.DefaultServiceRules */}}
{{- "# ServiceRules\n" -}}
{{- range $index, $service := $S -}}

	{{- $serviceGroups := tmpl.Exec "ParseSlice" (dict "dict" $service "key" "Groups") | data.JSONArray -}}
	{{- range $index, $group := $serviceGroups -}}
		{{- if not (has $Groups $group) -}}
			{{- $Groups = $Groups | append $group -}}
		{{- end -}}
	{{- end -}}

	{{- $name := index $service "Name" -}}
	{{/* Load service rules */}}
	{{- $SR := tmpl.Exec "ParseDict" (dict "dict" $service "key" "Rules") | data.JSON -}}

	{{- $ruleNames := keys $SR -}}
	{{- range $index, $rule := $MR -}}
		{{- if not (has $ruleNames $rule) -}}
			{{- $ruleNames = $ruleNames | append $rule -}}
		{{- end -}}
	{{- end -}}

    {{/* Render service rules */}}
    {{- range $index, $rule := $ruleNames -}}

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
{{- "\n" -}}

{{/* Render grouped rules */}}
{{- "# GroupedRules\n\n" -}}
{{- range $index, $group := $Groups -}}
	{{- if eq $group "_" -}}
		{{- "# Main Rules\n" -}}
	{{- else if eq $group "_all" -}}
		{{- continue -}}
	{{- else -}}
		{{- printf "# %s Rules\n" $group -}}
	{{- end -}}

	{{- range $index, $rule := $MR -}}
		{{- $dependencies := "" -}}
	
		{{- range $index, $service := $S -}}
			{{- $serviceGroups := tmpl.Exec "ParseSlice" (dict "dict" $service "key" "Groups") | data.JSONArray -}}

			{{- if or (or (eq $group "_") (has $serviceGroups $group)) (has $serviceGroups "_all") -}}
				{{- $dependencies = printf "%s %s_%s" $dependencies (index $service "Name") $rule -}}
			{{- end -}}
		{{- end -}}

		{{- $dependencies = strings.TrimPrefix " " $dependencies}}

		{{- $groupRule := $rule -}}
		{{- if ne $group "_" -}}
			{{- $groupRule = printf "%s_%s" $group $rule -}}
		{{- end -}}

		{{- template "RenderRuleWithoutNewLine" dict "rule" $groupRule "dependencies" $dependencies "command" "" -}}
	{{- end -}}

	{{ printf "\n" }}
{{- end -}}
`
