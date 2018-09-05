# firemodel [![CircleCI](https://circleci.com/gh/mickeyreiss/firemodel.svg?style=svg)](https://circleci.com/gh/mickeyreiss/firemodel)

> Type-safe, cross-platform models for [Firestore](https://firebase.google.com/docs/firestore/).

Google's Firestore offering us a fantastic way to add "real time" to your product in need minutes. For simple products, what they provide of the box makes a lot of sense. You get dynamic, reflection-based API wrappers that for in with your existing code.

This approach breaks down, however, as you start to build against a complex or quickly changing data model, or when multiple developers need to interact with the same schema across multiple codebases.

Firemodel generates consistent models for all of your codebase automatically. Based on a single source of truth, your schema, firemodel generates idiomatic models for your firestore documents. These models are designed to work nicely with the standard Firestore SDKs and its 7 native types.

Firemodel currently supports iOS (Swift+Pring), TypeScript and Go.

## Quick Start

Install firemodel with go get:

    go get github.com/mickeyreiss/firemodel/firemodel

Now that you have the tool, you're ready to define a schema!

_Note: this guide assumes you are already set up with Firestore in your project._

### 1. Define a Schema

Define a schema for all of the document types you work with. 

`schema.firemodel`:

```
// An Aircraft is a machine that flies.
model Aircraft {
  // The airplane's official registration. 
  string tailnumber;
  // The plane's current altitude. 
  integer altitude;
  // The airplane's current autopilot setting.
  AutoPilotMode ap;
}

enum AutoPilotMode {
  manual,
  heading,
  alt,
  fullAuto,
}
```

The Schema is a platform independent way to describe the data you want to store in Firestore.

You can define as many models and enums as you'd like. See the [example](firemodel.example.firemodel) for a complete description of the type system.

### 2. Generate models

Open up a terminal, and generate your models:

    firemodel compile --go_out=./gen/go --ts_out=./.gen/ts --ios_out=./.gen/ios --schema=schema.firemodel
    
This generated some Swift, some typescript and some go code. You'll find it in `.gen` directory, as requested. You can now incorporate these generated files into your project.

This is the standard firemodel workflow. Whenever you need to update your data model, you'll update the schema and regenerate the models. 

## 3. Use the models

The models are designed to be idiomatic for their target languages and the official Firestone SDKs. 

In go, firemodel provides you with a tagged struct.

In iOS, firemodel provides a [Pring](https://github.com/1amageek/Pring/) `Object` subclass.

In typescript, firemodel provides interfaces and helpers classes.

It is trivial to extend firemodel with custom language providers. See the `Modeler` interface for more details.

***

## Schema Language & Type System

With Firemodel, you use a language loosely inspired by proto3 to define your data model.

In firemodel, whitespace is generally ignored, and semi-colons are required.

A model has fields, and each field has a type: 

There's no type for nil, because it's silly to have a field that is always nil. Surprisingly, firestore represents nil as a type; that's not faithfully represented in firemodel.

The primitive types in firemodel match firestore's built-in types exactly:

| Firemodel Type  |  Firestore Type |
| --------------- | --------------- |
| `string`        | String          |
| `integer`       | Number          |
| `double`        | Double          |
| `timestamp`     | Timestamp       |
| `bytes`         | Bytes           |
| `reference`     | Reference       |
| `geopoint`      | GeoPoint        |
| `array`         | Array           |
| `map`           | Map             |

You can also define Enums:

```
enum TodoState {
  TODO,
  DONE,
}
```

These enum values can be used as a field type:

```
model Todo {
  TodoState status;
}
```

Enums are not real. They end up getting represented in each language, and in firestore itself as strings.

You can also embed a model type:

```
model Child {}

model Parent {
  Child daughter;
}
```

Embedded models are not real. They end up getting stored as a `Map` in firestore.

`collection` provides a nested collection. Collections are not real; they are not even fields. Collections do not get stored directly in firestore.

### Generics

Firemodel supports generics for `map`, `array` and `reference`.

Generic types can be firemodel primitives or user-defined types:

```
model Child {}

enum Emotion { HAPPY, SAD, }

model Thing {
  array<string> primitive_array;
  array<Child> model_array;
  array<Emotion> enum_array;
}
```

Generic generic types are not currently supported: `array<reference<T>>`.

### Options

You can specify schema and model options via the following syntax:

```
option foo.bar = "baz";

model MyModel {
    option lang.key = "value";
}
```

Options are used to provide hints to the modelers.

Here are the currently supported options:

| Option Name | Description | Example |
| --------- | ------------ | ---- |
| `firestore.path` | Document's typical location in firestore, specified as a template with variables surrounded with curly braces. | `option firestore.path = "users/{user_id}";` |
| `firestore.autotimestamp` | Automatically add createdAt and updatedAt fields. | `option firestore.autotimestamp = true;` |
| `ts.namespace` | The TypeScript namespace for generated interfaces. | `option ts.namespace = "SomeApp";` |
| `go.package` | The name of the go package for generated code. | `option go.package = "myapp";` |

 
