<!-- This comment is uncommented when auto-synced to www-kluctl.io

---
title: Templating
description: Templating documentation.
weight: 30
---
-->

# Templating

The Template Controller reuses the Jinja2 templating engine of [Kluctl](https://kluctl.io).

Documentation is available [here](https://kluctl.io/docs/kluctl/templating/).

## Predefined variables

You can use multiple predefined variables in your templates. These are:

### objectTemplate

Available in templates inside [ObjectTemplate](./spec/v1alpha1/objecttemplate.md) and represents the whole
`ObjectTemplate` that was on your target BEFORE the reconciliation started.

### textTemplate

Available in templates inside [TextTemplate](./spec/v1alpha1/texttemplate.md) and represents the whole
`TextTemplate` that was on your target BEFORE the reconciliation started.
