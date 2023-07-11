import {initializeApp} from 'https://www.gstatic.com/firebasejs/10.0.0/firebase-app.js'
import {getAuth, signInWithEmailAndPassword} from "https://www.gstatic.com/firebasejs/10.0.0/firebase-auth.js"

const firebaseApp = initializeApp({
    apiKey: "AIzaSyAbGRjnPPIOyHj0vUlMYb10V8d32OJpzxY",
    authDomain: "inventory-management-1f296.firebaseapp.com",
    projectId: "inventory-management-1f296",
    storageBucket: "inventory-management-1f296.appspot.com",
    messagingSenderId: "302262208094",
    appId: "1:302262208094:web:73ad1b32c5a63ea99316cf",
    measurementId: "G-Y3VZCX11RN"
})
const auth = getAuth()

const loginBtn = document.getElementById("login-btm")
loginBtn.addEventListener("click", function (e) {
    const email = document.getElementById("email").value
    const password = document.getElementById("password").value

    signInWithEmailAndPassword(auth, email, password)
        .then((userCredential) => {
            // Signed in
            const user = userCredential.user
            const firebaseIDToken = document.getElementById("firebase-id-token")

            firebaseIDToken.innerHTML = `Firebase ID token: ${user.accessToken}`
        })
        .catch((error) => {
            const errorCode = error.code
            const errorMessage = error.message

            const errorDiv = document.getElementById("error")
            errorDiv.innerHTML = `<div>Error Code: ${errorCode}</div><br><div>Error Message: ${errorMessage}</div>'`
        })
})