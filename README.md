> [!WARNING]
> Under construction!

# Overview

Welcome to the Skal language repository!

Skal is a syntactically simple Rust-y syntax which transpiles to Lua (5.1).

```
# Define Unit types.
pub enum UnitType {
  FRIEND = 'friend'
  FOE    = 'foe'
}

# Define a global Unit struct.
pub struct Unit {
  Name
  HP
  AttackPower
  Type

  # Constructor
  new(name, t, power, hp) {
    this.Name        = name
    this.Type        = t
    this.AttackPower = power
    this.HP          = hp
    return this
  }

  # Method
  heal(points) {
    this.HP = this.HP + points
  }

  # Method
  damage(points) {
    this.HP = this.HP - points
  }

  # Method
  attack(other_unit) {
    other_unit.HP = other_unit.HP - this.AttackPower
  }
}

# Instantiate a Unit.
friend = Unit('Champion', UnitType.FRIEND, 7, 36)
foe    = Unit('Undead', UnitType.FOE, 12, 21)

# Attack!
friend.attack(foe)
print(foe.HP)
```

# Why Another Transpile to Lua Language?

Lua is a beautifully simple language and there are many existing transpile-to-lua languages: Teal, Erde, ClueLang, Yuescript, Haxe and many more. However, when
reading through these language docs I couldn't escape these feelings of, "well..
I don't love that..", or, "it'd be great if it just had this..". Also, none of
them really solved the things with Lua I found the most frustrating.

I'm a big proponent of developer ergonomics and experience and I want a language
that's as syntactically light as possible without sacrificing capability.

## What Problems is Skal Intended to Solve?

### Readability

The largest is readability. While syntactically and semantically very simple,
the Lua language relies almost entirely on words rather than symbols. This
results in some incredbily "noisy" code. When you introduce things like
comment-based type annotations your code just becomes a wall of text and it
becomes difficult to hone in quickly on the bits you're looking for when jumping
around the codebase.

Some implementation examples:
- `&&` and `||` rather than `and` and `or`
- Scoping behavior is inverted: rather than declaring everything `local` and
omitting `local` when you want something to be `global`, in Skal you simply
declare something `pub` when you want it to be global.
- `{` and `}` rather than `then` and `end`.
- `!` rather than `not`.
- Shorter keywords in general, such as `fn` rather than `function`.
- Method definitions are contained _within_ the struct definition.

A basic concrete example:

Lua:
```lua
local my_obj = {}
my_obj.__index = my_obj
function my_obj:my_function(arg1, arg2)
  if not arg1 or not arg2 or arg1 > arg2 then
    return arg2
  end
end
```

Skal:
```
struct my_obj {
  my_function(arg1, arg2) {
    if !arg1 || !arg2 || arg1 > arg2 {
      return arg2
    }
  }
}
```

### A More Modern Language Feel

As mentioned earlier, the Skal syntax is very Rust-y. The idea is to provide a
feeling of a lovely modern syntax which will feel like a more natural transition
for natives of more modern languages like Go and Rust.

Some implementation examples:
- Global scope declarations are controlled via the `pub` keyword.
- We use `struct`s and `enum`s rather than `table`s (tables do not exist in Skal).
- Some modern trappings such as `defer` and arrow functions (lambdas) are supported.

### Type safety!

Lua's extremely basic and dynamic type system coheses well with it's general
mantra of simplicity. But sometimes, I'd just like to know ahead of time if I
mistakenly sent an integer to a function expecting a string where that string is
going to cause a terminating error.

Skal's type system is currently under development but the aim is a happy medium
between requiring explicit type hints and inferring types where it's particularly
obvious.

Some example code:

```
# These would all be inferred.
x = 12             # int
y = x              # int
a = 'abc'          # int
b = 'abc' .. 'def' # str
c = b .. a         # str
d = 1 + (2 / 12)   # int
e = d / (14 * d)   # int

# Function args and return types (if they return a value) always require
# explicit annotation.
#
# There are no plans to implement a full Hindley-Milner type system.
pub fn some(arg1: str, arg2: str) str {
  return arg1 .. arg2
}
```

# Language Feature Status

| Feature                                | Status | Notes |
| -------------------------------------- | ------ | ----- |
| Undefined Reference Detection          | ✔️      |       |
| Skal Standard Library                  | ♻️      |       |
| Type System                            | ❌      |       |
| Pattern Matching, Algebraic Data Types | ❌      |       |

# Tooling Support Status

| Feature                             | Status | Notes                                                           |
| ----------------------------------- | ------ | --------------------------------------------------------------- |
| Syntax Highlighting, Brace Matching | ✔️      |                                                                 |
| Language Server                     | ❌      |                                                                 |
| Linter                              | ❌      |                                                                 |
| Formatter                           | ❌      |                                                                 |
| Skal Interpreter                    | ✔️      | Initial release uses Gopher-Lua, homegrown interpreter to come. |
| Skal LLVM Backend Support           | ❌      |                                                                 |
