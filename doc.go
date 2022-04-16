// A simple helper package for better constants usage in go.
// The packapge is not simply `constants` as it's a very
// generic name leading to high risk of colision with others.
//
// The package contains a Metadata struct which implements some
// helpers. Create a Metadata variable for constants, and write
// small methods wrapping call to the helpers. Just pick helpers
// you need and ignore the others.
//
// The helpers are convenient but not the fastest. In a critical path
// consider using way to do the job, like using switch/case instead of maps.
//
// See example for usages.
package goconstants
