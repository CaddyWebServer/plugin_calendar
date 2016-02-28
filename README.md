# Calendar

**Note: This package is not yet complete**

This package adds calendar management functionality to Abot's training
interface, providing a way for human trainers to build .

Put `m.component(abot.Gcal)` into /assets/js/profile.js to enable users to add
and remove their Google calendar account from Abot.

Add the following to your route definitions (in handlers.go) to enable creating
an event via the training interface.

```
// OAuth routes
e.Post("/api/calendar/events.json", handlerAPICalendarEventCreate)
e.Post("/oauth/connect/gcal.json", handlerOAuthConnectGoogleCalendar)
e.Post("/oauth/disconnect/gcal.json", handlerOAuthDisconnectGoogleCalendar)
```

Also add to that same profile.js page the following to ensure that the required
Google libraries are loaded.

```
abot.loadJS("https://apis.google.com/js/client.js?onload", function() {
	gapi.load("auth2", ctrl.auth2Callback)
})

ctrl.auth2Callback = function() {
	abot.auth2 = gapi.auth2.getAuthInstance()
	if (!!abot.auth2) {
		return
	}
	var gid = document.querySelector("meta[name=google-client-id]").getAttribute("content")
	gapi.auth2.init({
		client_id: gid,
		scope: "https://www.googleapis.com/auth/calendar"
	}).then(function(a) {
		abot.auth2 = a
		if (abot.auth2.isSignedIn.get()) {
			var email = abot.auth2.currentUser.get().getBasicProfile().getEmail()
			ctrl.toggleGoogleAccount(email)
		}
	}, function(err) {
		console.error(err)
	})
}

ctrl.toggleGoogleAccount = function(name) {
	var googleLink = document.getElementById("oauth-google-success-a")
	if (!googleLink) {
		// Not on the Profile page. This function is called globally on
		// Google's script loading, so it isn't dependent on any route.
		// Ultimately Google's script should only load on the Profile route,
		// which eliminates the need for this check
		return
	}
	var signout = document.getElementById("oauth-google-success")
	var signin = document.getElementById("signinButton")
	if (!name) {
		googleLink.text = ""
		signout.classList.add("hidden")
		signin.classList.remove("hidden")
	} else {
		googleLink.text = "Google - " + name
		signout.classList.remove("hidden")
		signin.classList.add("hidden")
	}
}
```
