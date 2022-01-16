import * as React from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import { useRouter } from 'next/router'
import axios from 'axios';



export default function PassengerSignUp() {
  const router = useRouter()
  
    async function handleSubmit(event) {
      event.preventDefault();
      const data = new FormData(event.target);
      const jsonString = JSON.stringify({
        firstname: data.get("firstname"),
        lastname: data.get("lastname"),
        moblieno: data.get("moblieno"),
        emailaddress: data.get("email"),
      });
      console.log(jsonString);
      // const res = await fetch('http://ridely-passenger:5000/api/v1/passenger/createPassenger', {
      //   body : jsonString,
      //   method : 'POST',
      //   headers : {
      //     'Content-Type' : 'application/json',
      //   },
      //   mode : 'no-cors',
      // })
      // const response = await res.json();
      // console.log(response);
      // router.push('/');
      axios.post("http://passenger/api/v1/passenger/createPassenger", 
        jsonString)
      .then(function (response) {
        console.log()
        response.status === 200 ? router.push('/') : console.log(response.status);
      })
      .catch(function (error) {
        console.log(error);
      });
         
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