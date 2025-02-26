import { useQuery, UseQueryResult } from "react-query";
import { Organization } from "./types";
import { get } from "./api";

export default function useGetOrganization(disabled?: boolean): UseQueryResult<Organization> {
  return useQuery("organization-workspaces", () => getOrganization(), { enabled: !disabled });
}

async function getOrganization(): Promise<Organization> {
  const response = await get("organization-workspaces");

  if (!response.ok) throw response.error;

  return response.data as Organization;
}
