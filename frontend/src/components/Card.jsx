import { useState } from "react";
import { useTheme } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import AutoGraphIcon from "@mui/icons-material/AutoGraph";
import BarChartIcon from "@mui/icons-material/BarChart";
import CardModel from "./Model";
import { Paper } from "@mui/material";

export default function GraphCard({ count, name, color }) {
  const theme = useTheme();
  const [open, setOpen] = useState(false);

  return (
    <>
      <Paper onClick={() => setOpen(true)}>
        <Card
          sx={{
            display: "flex",
            width: "40vw",
            backgroundColor: `${color}`,
            justifyContent: "space-between",
            "&: hover": {
              cursor: "pointer",
              boxShadow: `${theme.shadows[5]}`,
            },
          }}>
          <Box sx={{ display: "flex", flexDirection: "column" }}>
            <CardContent sx={{ flex: "1 0 auto" }}>
              <Typography component="div" variant="h5" color={"white"}>
                {count}
              </Typography>
              <Typography
                variant="subtitle2"
                color="whitesmoke"
                component="div">
                {name}
              </Typography>
            </CardContent>
          </Box>
          {name === "No of consumers" ? (
            <AutoGraphIcon
              sx={{ color: "white", fontSize: "4rem", pr: "1rem", pt: "1rem" }}
            />
          ) : (
            <BarChartIcon
              sx={{ color: "white", fontSize: "4rem", pr: "1rem", pt: "1rem" }}
            />
          )}
        </Card>
      </Paper>
      <CardModel open={open} handleClose={() => setOpen(false)} name={name} />
    </>
  );
}
