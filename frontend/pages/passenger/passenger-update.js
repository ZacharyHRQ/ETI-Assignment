import * as React from 'react';
import axios from 'axios';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Select from '@mui/material/Select';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import { useRouter } from 'next/router'
import { useEffect } from 'react';


export async function getStaticProps() {
  const res = await axios.get('http://passenger:5000/api/v1/passengersid/')
  const passengersid = await res.data;
  return {
    props: { passengersid }
  }
}

export default function PassengerUpdate({passengersid}){
  const router = useRouter()
    const [id, setid] = React.useState("");
    const [firstname, setfirstname] = React.useState("");
    const [lastname, setlastname] = React.useState("");
    const [moblieno, setmoblieno] = React.useState("");
    const [emailaddress, setaddress] = React.useState("");

    const handleChange = (event) => {
      setid(event.target.value);
      console.log(id);
    }

    useEffect(() => {
      if (id !== "") {
        axios.get('http://localhost:5000/api/v1/passenger/'+id)
        .then(res => {
          const passenger = res.data;
          setfirstname(passenger.firstname);
          setlastname(passenger.lastname);
          setmoblieno(passenger.moblieno);
          setaddress(passenger.emailaddress);
        })
      }
    },[id]);

    async function handleSubmit(event) {
        event.preventDefault();
        const jsonString = JSON.stringify({
          firstname: firstname,
          lastname: lastname,
          moblieno: moblieno,
          emailaddress: emailaddress,
        });
        console.log(jsonString);
        const res = await fetch('http://localhost:5000/api/v1/passenger/updatePassenger/'+id, {
          body : jsonString,
          method : 'POST',
          headers : {
            'Content-Type' : 'application/json',
          },
          mode : 'no-cors',
        })
        console.log(res.status);
        if (res.status === 0){ 
          alert("Update Successfully");
          router.push('/')
        }
      }

    return (<div> 
    
    <Container component="main" maxWidth="xs">
      <h1>Passenger Update</h1>
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
              id="email"
              label="Email Address"
              name="emailaddress"
              onChange={(e) => setaddress(e.target.value)}
              value={emailaddress}
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="firstname"
              label="FirstName"
              name="firstname"
              onChange={(e) => setfirstname(e.target.value)}
              value={firstname}
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="lastname"
              label="LastName"
              name="lastname"
              onChange={(e) => setlastname(e.target.value)}
              value={lastname}
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="moblieno"
              label="Moblie Number"
              id="moblieno"
              onChange={(e) => setmoblieno(e.target.value)}
              value={moblieno}
              inputProps={{maxLength:8}}
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


