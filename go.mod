module github.com/porticolabs/pbd-t

go 1.15

// This replace is a temporal workaround to get media upload features
// Should be removed once this PR gets merged: https://github.com/dghubble/go-twitter/pull/148
replace github.com/dghubble/go-twitter => github.com/janisz/go-twitter v0.0.0-20201206102041-3fe237ed29f3

require (
	github.com/adjust/rmq/v3 v3.0.0
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.7.0
	github.com/sirupsen/logrus v1.8.0
)
