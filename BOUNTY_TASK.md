# Bounty Fix

**Reward:** $10.00
**Issue:** https://github.com/Kristywvs22/gin/issues/1
**Repo:** https://github.com/Kristywvs22/gin

---

## Issue

🎯 Fix Data Race in (*Context).Copy() by Shallow Copying the Keys Map

## 📝 Description
The `(*Context).Copy()` method is the recommended way to pass a Gin context to a read-only asynchronous goroutine. However, the `Keys` map (which stores metadata via `Context.Set()`) is copied by reference rather than value. 

Because Go maps are not thread-safe, if the original request goroutine continues to modify context keys (e.g., in subsequent middleware or handlers) while an asynchronous background goroutine reads or writes to the copied context's keys, a concurrent map access violation occurs. This leads to:
- Data races (detected via `go test -race`).
- Inconsistent values returned by `Context.Get()`.
- Fatal runtime panics: `fatal error: concurrent map read and map write` or `concurrent map writes`.

This issue heavily impacts logging, tracing, audit middleware, and background task dispatchers that rely on `c.Copy()` to safely offload work.

## 🎯 Acceptance Criteria
- [ ] `(*Context).Copy()` must allocate a new `map[string]any` for the copied context's `Keys` field if the original context has initialized keys.
- [ ] All key-value pairs from the original `Keys` map must be copied into the new map (shallow copy of the map).
- [ ] Modifying the `Keys` map in the copied context must not affect the `Keys` map in the original context, and vice versa.
- [ ] If `c.Keys` is `nil`, the copied context's `Keys` should also remain `nil` (or safely initialized to avoid unnecessary allocations if no keys exist).
- [ ] A regression test must be added to verify that concurrent writes/reads to the original and copied context keys do not trigger the Go race detector.

## 🛠️ Technical Specifications & Context
The issue lies within the `Copy()` method of the `Context` struct, typically located in `context.go`.

### Current Behavior (Conceptual)
```go
func (c *Context) Copy() *Context {
    cp := Context{
        writermem: c.writermem,
        Request:   c.Request,
        Params:    c.Params,
        Keys:      c.Keys, // Reference copy of the map pointer
        // ...
    }
    // ...
    return &cp
}
```

### Proposed Fix
Modify the `Copy()` method in `context.go` to safely copy the map:

```go
func (c *Context) Copy() *Context {
    cp := Context{
        writermem: c.writermem,
        Request:   c.Request,
        Params:    c.Params,
        // ...
    }
    
    // Safely copy the Keys map
    if c.Keys != nil {
        cp.Keys = make(map[string]any, len(c.Keys))
        for k, v := range c.Keys {
            cp.Keys[k] = v
        }
    }
    
    // ... rest of the copy logic
    return &cp
}
```

## 🧪 Verification & Testing

### Automated Test Case
Add a test case in `context_test.go` (e.g., `TestContextCopyKeysRace`) to reproduce and verify the fix:

```go
func TestContextCopyKeysRace(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("foo", "bar")

	cCopy := c.Copy()

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: Write to original context keys
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			c.Set(fmt.Sprintf("key-%d", i), i)
		}
	}()

	// Goroutine 2: Write to copied context keys
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			cCopy.Set(fmt.Sprintf("copy-key-%d", i), i)
		}
	}()

	wg.Wait()
}
```

### Verification Command
Run the tests with the race detector enabled:
```bash
go test -race -run TestContextCopyKeysRace
```
The test must pass without triggering any race conditions or panics.

---

[![Opire Bounty](https://img.shields.io/badge/bounty-Opire-emerald?style=for-the-badge)](https://github.com/Kristywvs22/gin)
<br/>
<hr/>

<details><summary>This repo is using Opire - what does it mean? 👇</summary><br/>💵 Everyone can add rewards for this issue commenting <code>/reward 100</code> (replace <code>100</code> with the amount).<br/>🕵️‍♂️ If someone starts working on this issue to earn the rewards, they can comment <code>/try</code> to let everyone know!<br/>🙌 And when they open the PR, they can comment <code>/claim #1</code> either in the PR description or in a PR's comment.<br/><br/>🪙 Also, everyone can tip any user commenting <code>/tip 20 @Kristywvs22</code> (replace <code>20</code> with the amount, and <code>@Kristywvs22</code> with the user to tip).<br/><br/>📖 If you want to learn more, check out our <a href="https://docs.opire.dev">documentation</a>.</details>

---

## Process

### 1. Orient yourself first
Before touching any code:
- Read the README and any relevant docs
- Run: find . -type f | grep -v node_modules | grep -v .git | grep -v dist | head -60
- Grep for the key terms from the issue title to find the relevant files
- Read the files around the bug — understand what's currently happening

### 2. Find and read the tests
- Look for test files near the relevant code (*.test.ts, *.spec.js, __tests__, etc.)
- Read them — they tell you what behavior is expected
- Understand what a passing fix looks like before writing a line

### 3. Implement the fix
- Make the MINIMAL change that fixes the issue — nothing more
- Do not refactor, rename, reorganize, or clean up unrelated code
- Do not add dependencies unless the issue explicitly requires it
- Do not modify config files (package.json, tsconfig.json, etc.) unless that IS the fix
- One logical change, one commit

### 4. Verify
- Run the test suite: check package.json scripts for "test", "lint", "typecheck"
  Common commands: npm test, yarn test, pnpm test, pytest, go test ./..., cargo test
- If TypeScript is present: npx tsc --noEmit
- Fix any test failures your change caused — if a test was wrong before, update it

### 5. Commit and document
- git add -A
- git commit -m "fix: 🎯 Fix Data Race in (*Context).Copy() by Shallow Copying the Keys Map"
- Write AGENT_SUMMARY.md with:
  - Root cause: what was wrong and why
  - Fix: what you changed and the reasoning
  - Verification: how to confirm it works
  - Concerns: anything a reviewer should know

---

## Stop conditions

If ANY of these are true, write "UNCERTAIN: [reason]" at the top of AGENT_SUMMARY.md, make NO code changes, and stop:
- The fix requires changing more than ~5 files
- You cannot determine what the correct behavior should be
- The issue requires a database migration, API key, or external service
- You find 3+ open PRs already attempting this fix (high competition)
- The repo is not in a language/framework you can reason about confidently

A clean UNCERTAIN is better than a bad PR.
