import * as React from 'react';
import axios from 'axios';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Select from '@mui/material/Select';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';


export async function getStaticProps() {
  const res = await axios.get('http://localhost:5000/api/v1/passengersid/')
  const passengersid = await res.data;
  return {
    props: { passengersid }
  }
}

export default function RequestTrip({passengersid}) {
    const [id, setid] = React.useState("");
    const [pickUpPostal, setpickUpPostal] = React.useState("");
    const [dropOffPostal, setdropOffPostal] = React.useState("");

    const handleChange = (event) => {
      setid(event.target.value);
      console.log(id);
    }

    async function handleSubmit(event) {
        event.preventDefault();
        const jsonString = JSON.stringify({
            pickuppostalcode : pickUpPostal,
            dropoffpostalcode : dropOffPostal,
        });
        console.log(jsonString);
        const res = await fetch('http://localhost:5002/api/v1/request/'+id, {
          body : jsonString,
          method : 'POST',
          headers : {
            'Content-Type' : 'application/json',

          },
          mode : 'no-cors',
        })
        console.log(res.status);
      }

    return (<div> 
    
    <Container component="main" maxWidth="xs">
      <h1>Request a Trip</h1>
      <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
      <InputLabel id="demo-simple-select-standard-label">Id</InputLabel>
      <Select
          labelId="demo-simple-select-standard-label"
          id="demo-simple-select-standard"
          value={id}
          onChange={handleChange}
          name="id"
          label="id"
        >
          <MenuItem value="">
            <em>None</em>
          </MenuItem>

          {passengersid.map(passengerId => (
            <MenuItem value={passengerId} key={passengerId}>{passengerId}</MenuItem>
          ))}
        </Select>
            <TextField
              margin="normal"
              required
              fullWidth
              name="pickuppostalcode"
              label="PickUp Postal Code"
              id="pickUpPostal"
              onChange={(e) => setpickUpPostal(e.target.value)}
              value={pickUpPostal}
              inputProps={{maxLength:6}}
            />

            <TextField
              margin="normal"
              required
              fullWidth
              name="dropoffpostalcode"
              label="Drop Off Postal Code"
              id="dropoffpostalcode"
              onChange={(e) => setdropOffPostal(e.target.value)}
              value={dropOffPostal}
              inputProps={{maxLength:6}}
            />
            
            
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              Update
            </Button>

            
    </Box>
    </Container>

    
    </div>)

}


