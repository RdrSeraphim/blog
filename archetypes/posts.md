---
date: {{ .Date }}
lastmod: {{ .Date }}
title: '{{ replace .File.ContentBaseName "-" " " | title }}'
draft: true
slug: {{ .File.ContentBaseName }}
t: []
summary: ''
---
