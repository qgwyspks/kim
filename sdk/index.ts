import { IMClient, sleep } from "./sdk";

// @ts-ignore
const main = async () => {
    let cli = new IMClient("ws://localhost:8000", "ccc");
    let { status } = await cli.login()
    console.log("client login return -- ", status)
}

main()