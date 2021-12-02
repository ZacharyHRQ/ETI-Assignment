import * as React from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import { useRouter } from 'next/router'


export default function DriverSignUp() {
  const router = useRouter()
  
  async function handleSubmit(event) {
    event.preventDefault();
    const data = new FormData(event.target);
    const jsonString = JSON.stringify({
      firstname: data.get("firstname"),
      lastname: data.get("lastname"),
      moblieno: data.get("moblieno"),
      emailaddress: data.get("email"),
      carlicenseno: data.get("carlicenseno"),
      identificationnumber: data.get("identificationnumber"),
    });
    console.log(jsonString);
    const res = await fetch('http://localhost:5001/api/v1/driver/createDriver', {
      body : jsonString,
      method : 'POST',
      headers : {
        'Content-Type' : 'application/json',
      },
      mode : 'no-cors',
    })
    console.log(await res.status);
    if (res.status === 0) {
      alert("Creation Successfully");
      router.push('/')
    }
        
  }

  return (<div> 
    
    <Container component="main" maxWidth="xs">
      <h1>Driver SignUp</h1>
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
            <TextField
              margin="normal"
              required
              fullWidth
              name="carlicenseno"
              label="Car License No"
              id="carlicense"
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="identificationnumber"
              label="Identification Number"
              id="identificationnumber"
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