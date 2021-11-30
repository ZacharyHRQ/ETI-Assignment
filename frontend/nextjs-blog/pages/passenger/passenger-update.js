import * as React from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Select from '@mui/material/Select';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';

export default function PassengerUpdate() {

    const [passenger, setPassenger] = React.useState(
      {
        id: '',
        firstname: '',
        lastname: '',
        moblieno: '',
        emailaddress: '',
      });

    const handleChange = (event) => {
      setPassenger({
        ...passenger,
        [event.target.name]: event.target.value,
      });
      console.log(passenger);
    }


    // parse json data to Object and update to fields
    const handleClick = (event) => {
      event.preventDefault();
        fetch('http://localhost:5000/api/v1/passenger/1')
              .then(response => response.json())
              .then(data => setPassenger(data));
    }
    

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
        const url = 'http://localhost:5000/api/v1/';
        const res = await fetch('http://localhost:5000/api/v1/passenger/createPassenger', {
          body : jsonString,
          method : 'POST',
          headers : {
            'Content-Type' : 'application/json',
          },
          mode : 'no-cors',
        });
  
        console.log(await res.text());
        
        
      }

    return (<div> 
    
    <Container component="main" maxWidth="xs">
      <h1>Passenger Update</h1>
      <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
      <InputLabel id="demo-simple-select-standard-label">Id</InputLabel>
      <Select
          labelId="demo-simple-select-standard-label"
          id="demo-simple-select-standard"
          value={passenger.id}
          onChange={handleChange}
          name="id"
          label="id"
        >
          <MenuItem value="">
            <em>None</em>
          </MenuItem>
          {/* get ids from passenger services */}
          <MenuItem value={1}>1</MenuItem>
          <MenuItem value={2}>2</MenuItem>
          <MenuItem value={3}>3</MenuItem>
        </Select>
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Email Address"
              name="emailaddress"
              value={passenger.emailaddress}
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="firstname"
              label="FirstName"
              name="firstname"
              value={passenger.firstname}
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="lastname"
              label="LastName"
              name="lastname"
              value={passenger.lastname}
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="moblieno"
              label="Moblie Number"
              id="moblieno"
              value={passenger.moblieno}
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

    <Button
              onClick={handleClick}
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              Fetch Data
            </Button>
    </Container>
    </div>)
}