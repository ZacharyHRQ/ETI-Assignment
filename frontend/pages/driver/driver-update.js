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
  const res = await axios.get('http://localhost:5001/api/v1/fetchAllIds')
  const driverids = await res.data;
  return {
    props: { driverids }
  }
}

export default function PassengerUpdate({driverids}) {
    const [id, setid] = React.useState("");
    const [firstname, setfirstname] = React.useState("");
    const [lastname, setlastname] = React.useState("");
    const [moblieno, setmoblieno] = React.useState("");
    const [emailaddress, setaddress] = React.useState("");
    const [carlicense, setcarlicense] = React.useState("");


    const handleChange = (event) => {
      setid(event.target.value);
      console.log(id);
    }


    // parse json data to Object and update to fields
    const handleClick = () => {
      axios.get('http://localhost:5001/api/v1/driver/'+id)
      .then(res => {
        const driver = res.data;
        console.log(driver);
        setfirstname(driver.firstname);
        setlastname(driver.lastname);
        setmoblieno(driver.moblieno);
        setaddress(driver.emailaddress);
        setcarlicense(driver.carlicenseno);
      })
    }

    async function handleSubmit(event) {
        event.preventDefault();
        const jsonString = JSON.stringify({
          firstname: firstname,
          lastname: lastname,
          moblieno: moblieno,
          emailaddress: emailaddress,
          carlicenseno: carlicense
        });
        console.log(jsonString);
        const res = await fetch('http://localhost:5001/api/v1/driver/updateDriver/'+id, {
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
      <h1>Driver Update</h1>
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

          {driverids.map(driverId => (
            <MenuItem value={driverId} key={driverId}>{driverId}</MenuItem>
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
            <TextField
              margin="normal"
              required
              fullWidth
              name="carlicenseno"
              label="Car License No"
              id="carlicense"
              onChange={(e) => setcarlicense(e.target.value)}
              value={carlicense}
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


