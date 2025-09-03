# Hash Table / Dictionary Literal Syntax Comparison

This table compares how various programming languages represent **hash tables**, **dictionaries**, or **associative arrays** using literal syntax.

| Language        | Data Structure Name      | Literal Syntax Example                                      | Notes                                                                 |
|-----------------|---------------------------|--------------------------------------------------------------|-----------------------------------------------------------------------|
| **Python**      | `dict`                    | `{"key1": "value1", "key2": "value2"}`                       | Keys can be any immutable type.                                       |
| **JavaScript**  | `Object` or `Map`         | `{"key1": "value1", "key2": "value2"}`                       | Only strings (or symbols) as object keys in literal form.             |
| **TypeScript**  | `Record` / `{ [key]: val }`| `{ key1: "value1", key2: "value2" }`                         | Allows typed dictionaries; syntax same as JS for object literals.     |
| **Ruby**        | `Hash`                    | `{ "key1" => "value1", "key2" => "value2" }` or `{ key1: "value1" }` | The latter uses symbols as keys.                                      |
| **Go**          | `map[K]V`                 | `map[string]string{"key1": "value1", "key2": "value2"}`      | Key type must be comparable.                                          |
| **Rust**        | `HashMap` (from `std::collections`) | `use std::collections::HashMap; let mut map = HashMap::from([("key1", "value1"), ("key2", "value2")]);` | No literal syntax; use constructor macros or `from`.                  |
| **Java**        | `HashMap<K,V>`            | `Map<String, String> map = Map.of("key1", "value1", "key2", "value2");` | Since Java 9+. Earlier versions used `put()` repeatedly.              |
| **C#**          | `Dictionary<K,V>`         | `var dict = new Dictionary<string, string>{{"key1", "value1"}, {"key2", "value2"}};` | Requires constructor initialization.                                 |
| **PHP**         | `array` or `associative array` | `["key1" => "value1", "key2" => "value2"]`                  | Arrays can be used as both indexed and associative.                   |
| **Swift**       | `Dictionary<Key, Value>`  | `["key1": "value1", "key2": "value2"]`                       | Type inferred or declared explicitly.                                 |
| **Kotlin**      | `Map<K, V>`               | `mapOf("key1" to "value1", "key2" to "value2")`              | `mutableMapOf()` for mutable dictionaries.                            |
| **Elixir**      | `Map`                     | `%{"key1" => "value1", "key2" => "value2"}`                  | Atom keys: `%{key1: "value1"}`.                                       |
| **Lua**         | `table`                   | `{ key1 = "value1", key2 = "value2" }` or `{ ["key1"] = "value1" }` | Uses tables for all compound structures.                              |
| **Perl**        | Hash                      | `my %hash = ("key1" => "value1", "key2" => "value2");`       | Uses `=>` and scalar sigils (`%`).                                    |
| **Haskell**     | `Data.Map`                | `fromList [("key1", "value1"), ("key2", "value2")]`          | No built-in literal syntax for maps.                                  |
| **Julia**       | `Dict`                    | `Dict("key1" => "value1", "key2" => "value2")`               | Supports arbitrary key/value types.                                   |
| **Dart**        | `Map<K, V>`               | `{"key1": "value1", "key2": "value2"}`                       | Similar to JavaScript but with strong typing if specified.            |
