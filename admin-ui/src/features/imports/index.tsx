import Input from "../../components/Input";
import Pagination from "../../components/Pagination";
import useSyncPagination from "../../components/Pagination/hooks/useSyncPagination";
import Table from "../../components/Table";
import { Import } from "../../api/types";
import useApiKey from "../../api/useApiKey";
import useGetImports from "../../api/useGetImports";
import useGetOrganization from "../../api/useGetOrganization";
import useFilter from "../../hooks/useFilter";
import collectionCountLabel from "../../utils/collectionCountLabel";
import { importsTable } from "./utils/importsTable";
import style from "./style/Imports.module.scss";
import { PiMagnifyingGlassBold } from "react-icons/pi";

export default function Imports() {
  const { data: organization } = useGetOrganization();

  const workspaceId = organization?.workspaces?.[0]?.id || "";

  const { data: apiKey } = useApiKey(workspaceId);

  const { isLoading, data: imports } = useGetImports(workspaceId);

  const dataWithCombinedProperty = (imports || []).map((item) => ({
    ...item,
    importerNameAndId: `${item.importer?.name}_${item.id}`,
  }));

  const { dataFiltered, setFilter } = useFilter<Import[]>(["importerNameAndId"], dataWithCombinedProperty || []);

  const itemsPerPage = 25;
  const { dataPage, page, paginate, totalItems } = useSyncPagination(dataFiltered as any, itemsPerPage);

  const tableData = importsTable(dataPage as any, apiKey);

  if (isLoading) return null;

  return (
    <div className={style.imports}>
      <div className="container">
        <div className={style.header}>
          <div className={style.title}>
            <h1>Data</h1>
            {!isLoading && (
              <small>
                {collectionCountLabel(
                  "No files uploaded",
                  "file uploaded",
                  "files uploaded",
                  imports?.length || 0,
                  "No",
                  `file${dataFiltered.length === 1 ? "" : "s"} found from`,
                  dataFiltered.length
                )}
              </small>
            )}
          </div>

          <div className={style.actions}>
            <Input
              icon={<PiMagnifyingGlassBold />}
              type="search"
              className={style.searchInput}
              placeholder="Search"
              onChange={(e: any) => setFilter(e.target.value)}
            />
          </div>
        </div>

        {imports && <Table data={tableData} background="zebra" columnWidths={["30%", "25%", "10%", "30%", "5%"]} />}

        {!!(totalItems && totalItems > itemsPerPage) && (
          <Pagination totalItems={totalItems} itemsPerPage={itemsPerPage} onPageChange={paginate} initialPage={page} />
        )}
      </div>
    </div>
  );
}
