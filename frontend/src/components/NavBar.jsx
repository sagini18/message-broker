import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import Button from "@mui/material/Button";
import LogoutIcon from '@mui/icons-material/Logout';
import { Avatar, Typography } from "@mui/material";

export default function NavBar() {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static" sx={{backgroundColor:"#00256C"}}>
        <Toolbar sx={{display:"flex", justifyContent:"space-between"}}>
          {/* <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: 2 }}>
            <img src={`${stream_sync}`} alt="boardo" width="50px" />
          </IconButton> */}
          <Box display={"flex"}>
            <Avatar
              alt="Avatar"
              src="https://www.befunky.com/images/wp/wp-2021-01-linkedin-profile-picture-after.jpg?auto=avif,webp&format=jpg&width=944"
            />
            <Typography variant="h6" component="div" sx={{ flexGrow: 1}} paddingLeft={2} paddingTop={0.5}>
              Sagini Navaratnam
            </Typography>
          </Box>
          <Button color="inherit"> <LogoutIcon/></Button>
        </Toolbar>
      </AppBar>
    </Box>
  );
}
