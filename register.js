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
