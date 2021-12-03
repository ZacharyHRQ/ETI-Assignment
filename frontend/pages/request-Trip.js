import * as React from 'react';
import axios from 'axios';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Select from '@mui/material/Select';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import { useEffect } from 'react';


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
    const [trips, setTrips] = React.useState([]);

    useEffect(() => {
      if(id){
        axios.get('http://localhost:5002/api/v1/trip/'+id)
        .then(res => {
            setTrips(res.data);
            })

      }
    }, [id]);


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
              Request
            </Button>

            
        </Box>
        <h1>Previous Trips</h1>
    </Container>

    <TableContainer component={Paper}>
                <Table sx={{ minWidth: 650 }} aria-label="simple table">
                    <TableHead>
                    <TableRow>
                        <TableCell>TripId</TableCell>
                        <TableCell align="right">PassengerId</TableCell>
                        <TableCell align="right">DriverId</TableCell>
                        <TableCell align="right">PickUpPostalCode</TableCell>
                        <TableCell align="right">DropOffPostalCode</TableCell>
                        <TableCell align="right">TripStatus</TableCell>
                        <TableCell align="right">DateOfTrip</TableCell>

                    </TableRow>
                    </TableHead>
                    <TableBody>
                    {trips ? trips.map((row) => (
                        <TableRow
                        key={row.tripid}
                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                        <TableCell component="th" scope="row">
                            {row.tripid}
                        </TableCell>
                        <TableCell align="right">{row.passengerid}</TableCell>
                        <TableCell align="right">{row.driverid}</TableCell>
                        <TableCell align="right">{row.pickuppostalcode}</TableCell>
                        <TableCell align="right">{row.dropoffpostalcode}</TableCell>
                        <TableCell align="right">{row.tripstatus}</TableCell>
                        <TableCell align="right">{row.dateoftrip}</TableCell>
                        </TableRow>
                    )) : null}
                    </TableBody>
                </Table>
            </TableContainer>
  
    </div>)

}


