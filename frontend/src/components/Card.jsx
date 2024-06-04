import { memo, useState } from "react";
import { useTheme } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import AutoGraphIcon from "@mui/icons-material/AutoGraph";
import BarChartIcon from "@mui/icons-material/BarChart";
import CardModel from "./Model";
import { Paper } from "@mui/material";

export default memo(function GraphCard({ dataset, name, color, title }) {
  const theme = useTheme();
  const [open, setOpen] = useState(false);

  return (
    <>
      <Paper onClick={() => setOpen(true)}>
        <Card
          sx={{
            display: "flex",
            width: "35vw",
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
                {dataset?.length > 0 ? dataset[dataset?.length - 1]?.count : 0}
              </Typography>
              <Typography
                variant="subtitle2"
                color="whitesmoke"
                component="div">
                {name}
              </Typography>
            </CardContent>
          </Box>
          {name === "No of channels" ? (
            <BarChartIcon
              sx={{ color: "white", fontSize: "4rem", pr: "1rem", pt: "1rem" }}
            />
          ) : (
            <AutoGraphIcon
              sx={{ color: "white", fontSize: "4rem", pr: "1rem", pt: "1rem" }}
            />
          )}
        </Card>
      </Paper>
      <CardModel
        open={open}
        handleClose={() => setOpen(false)}
        title={title}
        name={name}
        color={color}
        dataset={dataset}
      />
    </>
  );
});
