import * as React from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';


export default function PassengerSignUp() {
  
    async function handleSubmit(event) {
      event.preventDefault();
      const data = new FormData(event.target);
      console.log({
        firstname: data.get("firstname"),
        lastname: data.get("lastname"),
        moblieno: data.get("moblieno"),
        emailaddress: data.get("email"),
      });
      const url = 'http://localhost:5000/api/v1/';
      const res = await fetch('http://localhost:5000/api/v1/passenger/createPassenger', {
        body : JSON.stringify({
          FirstName: data.get("firstname"),
          LastName: data.get("lastname"),
          MoblieNo: data.get("moblieno"),
          EmailAddress: data.get("email"),
        }),
        method : 'POST',
        headers : {
          'Content-Type' : 'application/json',
          'Access-Control-Allow-Origin' : '*',
          'Access-Control-Allow-Methods' : 'POST, GET, OPTIONS, PUT, DELETE',
          'Access-Control-Allow-Headers' : 'Content-Type, Accept, X-Requested-With, remember-me',
        },
        mode : 'no-cors',
      });
      
      const json = await JSON.parse(res);
      console.log(json);

    }

    return (<div> 
    
    <Container component="main" maxWidth="xs">
      <h1>Passenger SignUp</h1>
      <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Email Address"
              name="email"
              autoComplete="email"
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="firstname"
              label="FirstName"
              name="firstname"
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="lastname"
              label="LastName"
              name="lastname"
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="moblieno"
              label="Moblie Number"
              id="moblieno"
              inputProps={{maxLength:8}}
            />
            
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              Sign Up
            </Button>
    </Box>
    </Container>
    </div>)
}