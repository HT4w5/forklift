# Patching
Forklift can apply patches to fetches profiles.

## Patch config
```json
{
    "tag": "patch1",
    "content": {
        "+foo": [
            "bar1",
            "bar2"
        ]
    }
}
```

## Patch syntax
Object in `content` will be applyed on the original profile with the following rules:

### Default
Default action is deep merge. Simple types are overwritten.

### Modifiers
Add modifiers to keys to specify merge behavior.

#### Override
|Modifier|Type|Description|
|:--:|:--|:--|
|`!`|suffix|Force overwrite|


Example:
```yaml
# Original
foo:
  bar1: 1
  bar2: 2

# Patch
foo!:
- bar1
- bar2

# Result
foo:
- bar1
- bar2
```

#### Arrays
|Modifier|Type|Description|
|:--:|:--|:--|
|`+`|prefix|Prepend array|
|`+`|suffix|Append array|

Example:
```yaml
# Original
foo:
- bar1
- bar2

# Patch
+foo:
- bar3
- bar4

# Result
foo:
- bar3
- bar4
- bar1
- bar2
```
```yaml
# Original
foo:
- bar1
- bar2

# Patch
foo+:
- bar3
- bar4

# Result
foo:
- bar1
- bar2
- bar3
- bar4
```

#### Escape
Wrap keys that start or end with a modifier in brackets (`<>`).

Example:
```yaml
# Original
foo+:
  bar1: 1
  bar2: 2

# Patch
<foo+>!:
- bar1
- bar2

# Result
foo+:
- bar1
- bar2
```
