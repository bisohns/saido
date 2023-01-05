import React, { useEffect } from "react";
import { Container } from "@mui/material";
import { useParams } from "react-router-dom";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";

import PageHeader from "common/PageHeader";
import LoadingContent from "../common/LoadingContent";
import ServerDetailServicesTabPanel from "./ServerDetailServicesTabPanel";
import {
  ServerGroupedByHostResponseType,
  ServerResponseType,
} from "./ServerType";
import AppHeader from "AppHeader";

export default function ServerDetail({
  servers,
  connectionStatus,
  setJsonMessage,
}: {
  servers: ServerGroupedByHostResponseType;
  connectionStatus: string;
  setJsonMessage: (arg0: any) => void;
}) {
  const { host } = useParams<{ host: string }>();

  const [tabIndex, setTabIndex] = React.useState<number>(0);

  const services = servers[host as string];

  useEffect(() => {
    setJsonMessage({ FilterBy: host });
  }, [host]);

  const handleChangeTabIndex = (
    event: React.SyntheticEvent,
    newValue: number
  ) => {
    setTabIndex(newValue);
  };

  return (
    <>
      <AppHeader />{" "}
      <Container>
        <PageHeader
          title={`${host}`}
          breadcrumbs={[{ name: "Servers", to: "/" }, { name: `${host}` }]}
        ></PageHeader>

        <LoadingContent
          loading={connectionStatus === "Connecting"}
          error={connectionStatus === "Closed"}
        >
          <>
            <Tabs
              value={tabIndex}
              onChange={handleChangeTabIndex}
              aria-label={`${host} Tab`}
              variant="scrollable"
              scrollButtons="auto"
            >
              {services
                ?.sort()
                ?.map((service: ServerResponseType, index: number) => (
                  <Tab label={service.Message.Name} key={index} />
                ))}
            </Tabs>

            {services
              ?.sort()
              ?.map((service: ServerResponseType, index: number) => {
                // if (host !== servicesGroupedByName[serverName].Host)
                //   return null;

                return (
                  <div key={index}>
                    {index === tabIndex && (
                      <ServerDetailServicesTabPanel
                        serverName={service.Message.Name}
                        serverData={service}
                      />
                    )}
                  </div>
                );
              })}
          </>
        </LoadingContent>
      </Container>
    </>
  );
}
