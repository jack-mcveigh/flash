# CLI Usage

## Adding a Card

```bash
flash add -t "title" -d "A desc." <group>
```

OR

```bash
flash add -t "group.title" -d "A desc."
```

## Updating a Card

```bash
flash update -t "title" -d "New desc." <group>
```

OR

```bash
flash add -t "group.title" -d "New desc."
```

## Removing a Card

```bash
flash remove -t "title" <group>
```

OR

```bash
flash remove -t "group.title"
```

## Getting a Group

```bash
flash get <group>
```
