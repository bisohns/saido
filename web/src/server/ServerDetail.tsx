import React, { useEffect } from "react";
import { Container } from "@mui/material";
import { useParams } from "react-router-dom";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";

import PageHeader from "common/PageHeader";
import LoadingContent from "../common/LoadingContent";
import ServerDetailServicesTabPanel from "./ServerDetailServicesTabPanel";
import {
  ServerGroupedByNameResponseType,
  ServerResponseType,
  ServerServiceNameType,
} from "./ServerType";
import AppHeader from "AppHeader";

export default function ServerDetail({
  servicesGroupedByName,
  connectionStatus,
  setJsonMessage,
}: {
  servicesGroupedByName: ServerGroupedByNameResponseType;
  connectionStatus: string;
  setJsonMessage: (arg0: any) => void;
}) {
  const { host } = useParams<{ host: string }>();
  console.log("servicesGroupedByName", servicesGroupedByName);

  const [tabIndex, setTabIndex] = React.useState<number>(0);

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
              {Object.keys(servicesGroupedByName)
                ?.sort()
                ?.map((serverName: string, index: number) => (
                  <Tab label={serverName} key={index} />
                ))}
            </Tabs>

            {Object.keys(servicesGroupedByName)
              ?.sort()
              ?.map((serverName: string, index: number) => {
                if (host !== servicesGroupedByName[serverName].Host)
                  return null;

                return (
                  <div key={index}>
                    {index === tabIndex && (
                      <ServerDetailServicesTabPanel
                        serverName={serverName as ServerServiceNameType}
                        serverData={
                          servicesGroupedByName[
                            serverName as ServerServiceNameType
                          ]?.data?.at(-1) as ServerResponseType
                        } // get the last object of service
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
