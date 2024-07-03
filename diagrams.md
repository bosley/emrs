```
                                                                                   
                                                                               
                                                                               
                                │             │                                
            SERVER INSTANCE     │  SECTORS    │      ASSETS                    
                                │             │                                
     ───────────────────────────┼─────────────┼──────────────────────          
          Garden Monitor Server │             │                                
                  │                                                            
                  │                            ┌─────►soil sensor              
                  │                            │                               
                  ├─────────────► Flower Bed A └────► temp sensor              
                  │                                                            
                  │                                                            
                  │                            ┌────► soil sensor              
                  ├────────────►  Flower Bed B │                               
                  │                            └────► temp sensor              
                  │                                                            
                  │                        ┌────────► light sensor             
                  └─────────────► Tomatoes │                                   
                                           └─────────►humidity                 
                                                                              
                                                                                                                                                                
       /Garden/Tomatoes/Light Sensor                                                                                                                            
       │                                                                                                                                                        
       │                                                                                                                                                        
       │                                                                                                                                                        
       │      ┌────────────────────────►   HTTPS Endpoint                                                                                                       
       │      │                              │                           ┌────► consumer[action-runner<action-name>] ────┐                                      
  ┌────▼──────┴─────────┐                    ▼                           │                                               │                                      
  │ Event               │                  Pub/Sub [NERV]                │                                               │                                      
  │   Light Measurement │                  │                             │ ┌───►consumer[action-runner<action-name>] ────┤                                      
  │   Encoded           │                  │                             │ │                                             │                                      
  └─────────────────────┘                  ├────► topic[signal-name]─────┘ │                                             │                                      
                                           │                               ├─┬─►consumer[action-runner<action-name>]─────┤                                      
                                           ├────► topic[signal-name]───────┘ │                                           │                                      
                                           │                                 │                                           │                                      
                                           └────► topic[signal-name] ────────┘                                           │                                      
                                                                                                                         │                                      
                                                                                                                         │                                      
                                                                                                                         │         ┌───────────────────┐        
                                                                                                                         │  ┌─────►│ Twilio            │        
                                                                                     ┌───────────────────────────┐       │  │      │                   │        
                                                                                     │                           │       │  │      └───────────────────┘        
                                                                                     │ EMRS Runtime API          │       │  │                                   
                                                                                     │                           │ ◄─────┘  │      ┌───────────────────┐        
                                                                                     │   - In-Mem K/V Store      │          │ ┌───►│ Gmail             │        
                                                                                     │                           ├──────────┘ │    │                   │        
                                                                                     │   - 3rd Party module      │            │    └───────────────────┘        
                                                                                     │     interfacing           │◄───────────┘                                 
                                                                                     │                           │                 ┌───────────────────┐        
                                                                                     │                           │                 │ IoT Device(s)     │        
                                                                                     │                           │◄───────────────►│                   │        
                                                                                     └───────────────────────────┘                 └───────────────────┘        
                                                                                                                                                                
                                                                                                                         


                                                                                                                         ```
