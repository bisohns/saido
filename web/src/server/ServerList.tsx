import { Card, CardActionArea, Grid, Typography } from "@mui/material";
import { Box, Container } from "@mui/system";
import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { ReactComponent as ServerIcon } from "assets/svg/server.svg";
import LoadingContent from "common/LoadingContent";
import ThemeConfig from "ThemeConfig";
import { ServerGroupedByHostResponseType } from "./ServerType";
import AppHeader from "AppHeader";

export default function ServerList({
  servers,
  connectionStatus,
  setJsonMessage,
}: {
  servers: ServerGroupedByHostResponseType;
  connectionStatus: string;
  setJsonMessage: (arg0: any) => void;
}) {
  const navigate = useNavigate();

  useEffect(() => {
    setJsonMessage({ FilterBy: "" });
  }, []);

  return (
    <>
      <AppHeader />
      <Container>
        <LoadingContent
          loading={connectionStatus === "Connecting"}
          error={connectionStatus === "Closed"}
        >
          <Grid container spacing={2} my={10}>
            {Object.keys(servers)?.map((serverHost: string, index: number) => (
              <Grid item xs={12} md={6} key={index}>
                <Card
                  key={index}
                  style={{ background: ThemeConfig.palette.primary.light }}
                >
                  <CardActionArea onClick={() => navigate(`/${serverHost}`)}>
                    <Box
                      display={"flex"}
                      justifyContent="center"
                      alignItems="center"
                      flexDirection="column"
                      my={5}
                    >
                      <ServerIcon width={"100px"} />
                      <Typography
                        textTransform={"capitalize"}
                        mb={2}
                        noWrap
                        fontWeight={600}
                        style={{
                          color: ThemeConfig.palette.common.white,
                        }}
                      >
                        <>
                          {serverHost} -{" "}
                          <span
                            style={{
                              color: ThemeConfig.palette.success.dark,
                            }}
                          >
                            {servers[serverHost]?.[0]?.Message?.Platform}
                          </span>
                        </>
                      </Typography>
                    </Box>
                  </CardActionArea>
                </Card>
              </Grid>
            ))}
          </Grid>
        </LoadingContent>
      </Container>
    </>
  );
}
