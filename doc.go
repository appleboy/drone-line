// Sending line notifications using a binary, docker or Drone CI written in Go (Golang).
//
// Details about the drone-line project are found in github page:
//
//     https://github.com/appleboy/drone-line
//
// The pre-compiled binaries can be downloaded from release page.
//
// Setup Webhook service
//
// Setup Webhook service as default port 8088.
//
//   drone-line-v1.4.0-windows-amd64.exe \
//     --secret xxxx \
//     --token xxxx \
//     webhook
//
// Change default webhook port to 8089.
//
//   drone-line-v1.4.0-windows-amd64.exe \
//     --port 8089 \
//     --secret xxxx \
//     --token xxxx \
//     webhook
//
// Send Notification
//
// Setup the --to flag after fetch user id from webhook service.
//
//   drone-line-v1.4.0-windows-amd64.exe \
//     --secret xxxx \
//     --token xxxx \
//     --to xxxx \
//     --message "Test Message"
//
// For more details, see the documentation and example.
//
package main
