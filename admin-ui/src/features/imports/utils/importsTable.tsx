import { TableData } from "../../../components/Table/types";
import Download from "../components/Download";
import { Import } from "../../../api/types";
import { timeToText } from "../../../utils/time";

export function importsTable(imports: Import[] = [], apiKey: string): TableData {
  return imports.map((importItem) => {
    return {
      "Import ID": importItem.id,
      Importer: importItem.importer?.name,
      Rows: importItem.num_rows,
      Created: timeToText(importItem.created_at),
      _actions: {
        raw: importItem.id,
        content: <Download importItem={importItem} apiKey={apiKey} />,
      },
    };
  });
}
