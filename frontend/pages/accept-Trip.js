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



export async function getStaticProps() {
  const res = await axios.get('http://localhost:5001/api/v1/fetchAllIds')
  const driversid = await res.data;
  return {
    props: { driversid }
  }
}

export default function AcceptTrip({driversid}) {
    const [id, setid] = React.useState("");
    const [trips, setTrips] = React.useState([]);

    const handleChange = (event) => {
        setid(event.target.value);
    }

    const fetchTrips = () => {
        axios.get('http://localhost:5002/api/v1/getCurrentTrips/'+id)
        .then(res => {
            setTrips(res.data);
            })
    }

    async function acceptTrip(){
        const data = JSON.stringify({
            tripstatus : 1
        });
        const res = await fetch('http://localhost:5002/api/v1/changeStatus/'+id, { 
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: data,
            mode : 'no-cors',
        });
        console.log(res.status);
    
    }
        

    return (

        <Container maxWidth="m">
            <h1>Accept Trip</h1>

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

                {driversid.map(driverId => (
                    <MenuItem value={driverId} key={driverId}>{driverId}</MenuItem>
                ))}
            </Select>

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
                        <TableCell align="right">Accept</TableCell>

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
                        <Button
                            onClick={acceptTrip}  sx={{ mt: 3, mb: 2 }}>
                            Accept
                        </Button>

                        </TableRow>
                    )) : null}
                    </TableBody>
                </Table>
            </TableContainer>

            <Button
                onClick={fetchTrips}
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                >
                Update
            </Button>
        </Container>

        
            





    )

}
