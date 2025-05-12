const API_URL = "http://localhost:8080/todo/";

document.getElementById('register-form').addEventListener('submit', async function(event) {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirm-password').value;

    if (password !== confirmPassword) {
        alert('Passwords do not match');
        return;
    }

    try {
        const response = await fetch(`${API_URL}register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: email,
                password: password
            })
        });

        const data = await response.json();

        if (response.ok) {
            alert('Registration successful');
            window.location.href = 'login.html'; // Redirect to login page after successful registration
        } else {
            alert('Registration failed: ' + data.message);
        }
    } catch (error) {
        console.error('Error during registration:', error);
        alert('Error during registration: ' + error);
    }
});

// Example of registration
async function registerUser(email, password) {
    const response = await fetch('http://localhost:8080/todo/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            email: email,
            password: password
        })
    });
    
    if (response.ok) {
        const result = await response.json();
        console.log('User registered successfully:', result);
    } else {
        const error = await response.text();
        console.error('Error registering user:', error);
    }
}
