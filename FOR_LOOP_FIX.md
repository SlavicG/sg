# For-Loop Scoping Fix

## Problem Description

The original implementation had a scope inconsistency in for-loops:

### Original Buggy Code Flow:
```
for(let i = 0; i < 5; i = i + 1) {
    puts(i)
}
```

**What was happening:**
1. `let i = 0` → executed in **parent scope**
2. `i < 5` → evaluated in **iterationScope** (child of parent)
3. `puts(i)` → executed in **iterationScope** 
4. `i = i + 1` → executed in **parent scope** ❌ BUG!

**The Issue:**
- Each iteration created a fresh `iterationScope`
- The post-increment ran in the parent scope
- This caused scope conflicts and unpredictable behavior
- Variables declared in the initializer couldn't be properly updated

## Solution

Create a dedicated `loopScope` that wraps the entire for-loop:

### Fixed Code Flow:
```
for(let i = 0; i < 5; i = i + 1) {
    puts(i)
}
```

**What happens now:**
1. Create `loopScope` (child of parent scope)
2. `let i = 0` → executed in **loopScope** ✓
3. `i < 5` → evaluated in **loopScope** ✓
4. `puts(i)` → executed in **loopScope** ✓
5. `i = i + 1` → executed in **loopScope** ✓

**Benefits:**
- All parts of the loop share the same scope
- Loop variables are properly initialized and updated
- Nested loops work correctly (each has its own loopScope)
- No more workarounds needed (like using arrays `[0]`)

## Code Changes

**File:** `sg_interpreter/src/sg/evaluator/evaluator.go`

**Function:** `evalForStatement`

**Key Changes:**
1. Create `loopScope := Item.NewEnclosedScope(scope)` at the start
2. Execute initializer in `loopScope` instead of `scope`
3. Execute condition in `loopScope` instead of `iterationScope`
4. Execute body in `loopScope` instead of `iterationScope`
5. Execute post in `loopScope` instead of `scope`

## Test Cases

### Test 1: Simple Counter
```sg
for(let i = 0; i < 5; i = i + 1) {
    puts(i)
}
```
**Expected Output:**
```
0
1
2
3
4
```

### Test 2: Nested Loops
```sg
for(let i = 1; i <= 3; i = i + 1) {
    for(let j = 1; j <= 3; j = j + 1) {
        puts(i, j)
    }
}
```
**Expected Output:**
```
1 1
1 2
1 3
2 1
2 2
2 3
3 1
3 2
3 3
```

### Test 3: Modifying Outer Scope Variable
```sg
let sum = 0
for(let i = 1; i <= 5; i = i + 1) {
    sum = sum + i
}
puts("Sum:", sum)
```
**Expected Output:**
```
Sum: 15
```

### Test 4: Accessing Outer Scope
```sg
let multiplier = 2
for(let i = 1; i <= 3; i = i + 1) {
    puts(i * multiplier)
}
```
**Expected Output:**
```
2
4
6
```

## How to Test

1. Compile the interpreter:
   ```bash
   cd sg_interpreter/src/sg/main
   go build -o sg
   ```

2. Run the test file:
   ```bash
   ./sg ../../../test_for_loop.sg
   ```

3. Or test manually in REPL:
   ```bash
   ./sg
   ```

## Verification

After applying this fix:
- ✅ Loop counters work without array workarounds
- ✅ Nested loops maintain separate loop scopes
- ✅ Outer scope variables remain accessible
- ✅ Post-increment properly updates loop variable
- ✅ All examples from README should work correctly

## Impact

This fix resolves the critical scoping issue mentioned in the README TODO list:
> "Currently, the for loops don't work ideally and there are some problems with the scopes. It will be solved in a later version, where I will change the structuring of scopes and use a Scope-Stack."

The solution doesn't require a full scope-stack implementation - just ensuring all loop components share the same dedicated scope.
